import argparse
import asyncio
import logging
import signal
import sys
from pathlib import Path

import nats
from nats.errors import NoServersError

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

_NATS_CONNECT_RETRIES = 10
_NATS_CONNECT_DELAY = 3.0


async def _nats_error_cb(exc) -> None:
    logger.error("NATS error: %s", exc)


async def _nats_disconnected_cb() -> None:
    logger.warning("NATS disconnected")


async def _nats_reconnected_cb() -> None:
    logger.info("NATS reconnected")


async def _nats_closed_cb() -> None:
    logger.info("NATS connection closed")


async def _connect_nats(url: str) -> nats.aio.client.Client:
    for attempt in range(1, _NATS_CONNECT_RETRIES + 1):
        try:
            nc = await nats.connect(
                url,
                error_cb=_nats_error_cb,
                disconnected_cb=_nats_disconnected_cb,
                reconnected_cb=_nats_reconnected_cb,
                closed_cb=_nats_closed_cb,
            )
            logger.info("Connected to NATS at %s", url)
            return nc
        except NoServersError:
            if attempt == _NATS_CONNECT_RETRIES:
                raise
            logger.warning(
                "NATS not reachable at %s (attempt %d/%d), retrying in %.0fs...",
                url, attempt, _NATS_CONNECT_RETRIES, _NATS_CONNECT_DELAY,
            )
            await asyncio.sleep(_NATS_CONNECT_DELAY)


async def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("-config", default="config/config.yaml")
    args = parser.parse_args()

    cfg = load_config(Path(args.config))

    nc = await _connect_nats(cfg.nats.url)

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
