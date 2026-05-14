from minio import Minio

from document_analyzer.config import MinioConfig


class MinioStorage:
    def __init__(self, cfg: MinioConfig) -> None:
        self._client = Minio(
            cfg.endpoint,
            access_key=cfg.access_key,
            secret_key=cfg.secret_key,
            secure=cfg.secure,
        )
        self._bucket = cfg.bucket

    def get_object_bytes(self, object_key: str) -> bytes:
        response = self._client.get_object(self._bucket, object_key)
        try:
            return response.read()
        finally:
            response.close()
            response.release_conn()
