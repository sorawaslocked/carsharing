import asyncio
import logging

import grpc
from google.protobuf import empty_pb2

from service.user import document_pb2
from service.user import document_pb2_grpc
from document_analyzer.analyzer import DocumentAnalyzer
from document_analyzer.publisher import EventPublisher
from document_analyzer.storage import MinioStorage

logger = logging.getLogger(__name__)


class DocumentAnalyzerServicer(document_pb2_grpc.DocumentAnalyzerServiceServicer):
    def __init__(
        self,
        storage: MinioStorage,
        analyzer: DocumentAnalyzer,
        publisher: EventPublisher,
    ) -> None:
        self._storage = storage
        self._analyzer = analyzer
        self._publisher = publisher

    async def Analyze(
        self,
        request: document_pb2.AnalyzeRequest,
        context: grpc.aio.ServicerContext,
    ) -> empty_pb2.Empty:
        logger.info("Analyze request: document_id=%s object_key=%s", request.document_id, request.object_key)

        try:
            image_bytes = await asyncio.to_thread(self._storage.get_object_bytes, request.object_key)
        except Exception as exc:
            logger.error("MinIO fetch failed for %s: %s", request.document_id, exc)
            await context.abort(grpc.StatusCode.NOT_FOUND, f"Could not retrieve document: {exc}")
            return empty_pb2.Empty()

        try:
            # doctr is CPU-bound and synchronous — run it off the event loop
            result = await asyncio.to_thread(self._analyzer.analyze, image_bytes)
        except Exception as exc:
            logger.error("Analysis failed for %s: %s", request.document_id, exc)
            await context.abort(grpc.StatusCode.INTERNAL, f"Analysis error: {exc}")
            return empty_pb2.Empty()

        logger.info(
            "Document %s result: passed=%s defects=%s",
            request.document_id, result.passed, [d.type for d in result.defects],
        )

        try:
            await self._publisher.publish_document_analyzed(request.document_id, result)
        except Exception as exc:
            logger.error("Event publish failed for %s: %s", request.document_id, exc)
            await context.abort(grpc.StatusCode.INTERNAL, f"Event publish error: {exc}")
            return empty_pb2.Empty()

        return empty_pb2.Empty()
