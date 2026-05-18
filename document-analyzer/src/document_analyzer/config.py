from pathlib import Path
from typing import Any

import yaml
from pydantic import BaseModel


class GrpcConfig(BaseModel):
    host: str = "0.0.0.0"
    port: int = 50051
    max_workers: int = 10


class MinioConfig(BaseModel):
    endpoint: str = "localhost:9000"
    access_key: str = "minioadmin"
    secret_key: str = "minioadmin"
    secure: bool = False
    bucket: str = "documents"


class NatsConfig(BaseModel):
    url: str = "nats://localhost:4222"
    subject: str = "document.analyzed"


class AnalyzerConfig(BaseModel):
    min_width: int = 800
    min_height: int = 600
    blur_threshold: float = 80.0
    ocr_confidence_threshold: float = 0.5
    min_word_count: int = 3


class Config(BaseModel):
    grpc: GrpcConfig = GrpcConfig()
    minio: MinioConfig = MinioConfig()
    nats: NatsConfig = NatsConfig()
    analyzer: AnalyzerConfig = AnalyzerConfig()


def load_config(path: str | Path = "config/config.yaml") -> Config:
    path = Path(path)
    if not path.exists():
        return Config()
    with path.open() as f:
        data: dict[str, Any] = yaml.safe_load(f) or {}
    return Config(**data)
