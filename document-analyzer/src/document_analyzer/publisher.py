import logging

from nats.aio.client import Client as NatsClient

from event.user import document_pb2 as event_pb2
from document_analyzer.config import NatsConfig
from document_analyzer.analyzer import AnalysisResult

logger = logging.getLogger(__name__)


class EventPublisher:
    def __init__(self, nc: NatsClient, cfg: NatsConfig) -> None:
        self._nc = nc
        self._subject = cfg.subject

    async def publish_document_analyzed(self, document_id: str, user_id: str, result: AnalysisResult) -> None:
        event = event_pb2.DocumentAnalyzedEvent(
            document_id=document_id,
            user_id=user_id,
            passed=result.passed,
            defects=[
                event_pb2.Defect(type=d.type, description=d.description)
                for d in result.defects
            ],
        )
        await self._nc.publish(self._subject, event.SerializeToString())
        logger.info(
            "Published DocumentAnalyzedEvent for document %s to subject '%s' (passed=%s)",
            document_id, self._subject, result.passed,
        )
