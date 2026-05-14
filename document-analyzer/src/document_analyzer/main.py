import asyncio
import logging
import signal
import sys
from pathlib import Path

import nats

from document_analyzer import server
from document_analyzer.analyzer import DocumentAnalyzer
from document_analyzer.config import load_config
from document_analyzer.publisher import EventPublisher
from document_analyzer.servicer import DocumentAnalyzerServicer
from document_analyzer.storage import MinioStorage

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(name)s — %(message)s",
    stream=sys.stdout,
)
logger = logging.getLogger(__name__)


async def main() -> None:
    cfg = load_config(Path("config/config.yaml"))

    nc = await nats.connect(cfg.nats.url)
    logger.info("Connected to NATS at %s", cfg.nats.url)

    storage = MinioStorage(cfg.minio)
    analyzer = DocumentAnalyzer(cfg.analyzer)
    publisher = EventPublisher(nc, cfg.nats)
    servicer = DocumentAnalyzerServicer(storage, analyzer, publisher)

    grpc_server = server.create(cfg.grpc, servicer)
    await grpc_server.start()
    logger.info("gRPC server listening on %s:%d", cfg.grpc.host, cfg.grpc.port)

    # Intercept SIGINT/SIGTERM before asyncio can cancel the task.
    # shutdown.set() is called from the signal handler via call_soon_threadsafe,
    # so the await below returns normally instead of raising CancelledError.
    shutdown = asyncio.Event()
    loop = asyncio.get_running_loop()

    def _on_signal(sig: int, _frame) -> None:
        logger.info("Received %s, initiating shutdown...", signal.Signals(sig).name)
        loop.call_soon_threadsafe(shutdown.set)

    signal.signal(signal.SIGINT, _on_signal)
    signal.signal(signal.SIGTERM, _on_signal)

    await shutdown.wait()

    logger.info("Stopping gRPC server (5 s grace)...")
    await grpc_server.stop(grace=5)
    logger.info("Draining NATS connection...")
    await nc.drain()
    logger.info("Shutdown complete.")


def run() -> None:
    asyncio.run(main())


if __name__ == "__main__":
    run()
