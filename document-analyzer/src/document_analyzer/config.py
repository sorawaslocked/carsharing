import os
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


def _parse_bool(value: str) -> bool:
    return value.lower() in ("1", "true", "yes")


def load_config(path: str | Path = "config/config.yaml") -> Config:
    path = Path(path)
    data: dict[str, Any] = {}
    if path.exists():
        with path.open() as f:
            data = yaml.safe_load(f) or {}

    minio: dict[str, Any] = data.get("minio", {})
    if (v := os.getenv("MINIO_ENDPOINT")) is not None:
        minio["endpoint"] = v
    if (v := os.getenv("MINIO_USE_SSL")) is not None:
        minio["secure"] = _parse_bool(v)
    if (v := os.getenv("MINIO_ACCESS_KEY_ID")) is not None:
        minio["access_key"] = v
    if (v := os.getenv("MINIO_SECRET_ACCESS_KEY")) is not None:
        minio["secret_key"] = v
    if minio:
        data["minio"] = minio

    return Config(**data)
