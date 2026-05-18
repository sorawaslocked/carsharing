import grpc.aio

from service.user import document_pb2_grpc
from document_analyzer.config import GrpcConfig
from document_analyzer.servicer import DocumentAnalyzerServicer


def create(cfg: GrpcConfig, servicer: DocumentAnalyzerServicer) -> grpc.aio.Server:
    srv = grpc.aio.server()
    document_pb2_grpc.add_DocumentAnalyzerServiceServicer_to_server(servicer, srv)
    srv.add_insecure_port(f"{cfg.host}:{cfg.port}")
    return srv
