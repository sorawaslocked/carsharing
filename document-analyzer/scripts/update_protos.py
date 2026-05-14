"""
Re-download the vendored Python proto stubs from the canonical proto repo.
Run this whenever the proto repo tag is bumped.
"""
import urllib.request
from pathlib import Path

TAG = "v0.0.25"
BASE_URL = f"https://raw.githubusercontent.com/sorawaslocked/car-rental-protos/{TAG}/gen"
ROOT = Path(__file__).parent.parent / "src"

FILES = [
    "service/user/document_pb2.py",
    "service/user/document_pb2_grpc.py",
    "event/user/document_pb2.py",
    "event/user/document_pb2_grpc.py",
]

for rel in FILES:
    url = f"{BASE_URL}/{rel}"
    dest = ROOT / rel
    dest.parent.mkdir(parents=True, exist_ok=True)
    print(f"Fetching {url}")
    urllib.request.urlretrieve(url, dest)

# Ensure __init__.py files exist
for pkg in ["service", "service/user", "event", "event/user"]:
    (ROOT / pkg / "__init__.py").touch()

print("Done. Vendored stubs updated.")
