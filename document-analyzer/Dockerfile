FROM python:3.12-slim AS builder

WORKDIR /app

COPY pyproject.toml ./
COPY src ./src

# Install CPU-only PyTorch first to avoid downloading ~1 GB of CUDA packages.
# PYTHONPATH lets the second pip invocation see torch as already satisfied.
RUN PYVER=$(python -c "import sys; print(f'{sys.version_info.major}.{sys.version_info.minor}')") \
    && pip install --no-cache-dir --prefix=/install \
        --index-url https://download.pytorch.org/whl/cpu \
        torch torchvision \
    && PYTHONPATH=/install/lib/python${PYVER}/site-packages \
       pip install --no-cache-dir --prefix=/install .


FROM python:3.12-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    libglib2.0-0 \
    libgomp1 \
    libxcb1 \
    libxext6 \
    libx11-6 \
    libgl1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /install /usr/local
COPY config ./config

EXPOSE 9002

ENTRYPOINT ["document-analyzer", "-config", "./config/docker.yml"]
