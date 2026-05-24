import grpc.aio
from grpc_health.v1 import health, health_pb2, health_pb2_grpc

from service.user import document_pb2_grpc
from document_analyzer.config import GrpcConfig
from document_analyzer.servicer import DocumentAnalyzerServicer


def create(cfg: GrpcConfig, servicer: DocumentAnalyzerServicer) -> grpc.aio.Server:
    srv = grpc.aio.server()
    document_pb2_grpc.add_DocumentAnalyzerServiceServicer_to_server(servicer, srv)

    health_servicer = health.HealthServicer()
    health_pb2_grpc.add_HealthServicer_to_server(health_servicer, srv)
    health_servicer.set("", health_pb2.HealthCheckResponse.SERVING)

    srv.add_insecure_port(f"{cfg.host}:{cfg.port}")
    return srv
