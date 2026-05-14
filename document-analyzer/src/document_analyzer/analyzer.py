import io
import logging
from dataclasses import dataclass, field

import cv2
import numpy as np
from PIL import Image
from doctr.io import DocumentFile
from doctr.models import ocr_predictor

from document_analyzer.config import AnalyzerConfig

logger = logging.getLogger(__name__)

DEFECT_TOO_SMALL = "DOCUMENT_TOO_SMALL"
DEFECT_BLURRED = "BLURRED_IMAGE"
DEFECT_NO_CLEAR_TEXT = "NO_CLEAR_TEXT"


@dataclass
class Defect:
    type: str
    description: str


@dataclass
class AnalysisResult:
    passed: bool
    defects: list[Defect] = field(default_factory=list)


class DocumentAnalyzer:
    def __init__(self, cfg: AnalyzerConfig) -> None:
        self._cfg = cfg
        logger.info("Loading OCR predictor...")
        self._predictor = ocr_predictor(pretrained=True)
        logger.info("OCR predictor ready.")

    def analyze(self, image_bytes: bytes) -> AnalysisResult:
        image = Image.open(io.BytesIO(image_bytes)).convert("RGB")
        image_np = np.array(image, dtype=np.uint8)

        defects: list[Defect] = []

        # Check 1: image too small
        h, w = image_np.shape[:2]
        if w < self._cfg.min_width or h < self._cfg.min_height:
            defects.append(Defect(
                type=DEFECT_TOO_SMALL,
                description=(
                    f"Image dimensions {w}x{h} are below the minimum "
                    f"{self._cfg.min_width}x{self._cfg.min_height} required."
                ),
            ))

        # Check 2: blurred image via Laplacian variance
        gray = cv2.cvtColor(image_np, cv2.COLOR_RGB2GRAY)
        blur_score = float(cv2.Laplacian(gray, cv2.CV_64F).var())
        logger.debug("Blur score: %.2f (threshold: %.2f)", blur_score, self._cfg.blur_threshold)
        if blur_score < self._cfg.blur_threshold:
            defects.append(Defect(
                type=DEFECT_BLURRED,
                description=(
                    f"Blur score {blur_score:.2f} is below threshold "
                    f"{self._cfg.blur_threshold:.2f}. Image appears blurred."
                ),
            ))

        # Check 3: no clear text via OCR confidence
        doc = DocumentFile.from_images([image_np])
        result = self._predictor(doc)

        words = [
            word
            for page in result.pages
            for block in page.blocks
            for line in block.lines
            for word in line.words
        ]

        if len(words) < self._cfg.min_word_count:
            defects.append(Defect(
                type=DEFECT_NO_CLEAR_TEXT,
                description=(
                    f"Only {len(words)} words detected; minimum required is "
                    f"{self._cfg.min_word_count}. Document may lack legible text."
                ),
            ))
        elif words:
            avg_confidence = sum(w.confidence for w in words) / len(words)
            logger.debug("Avg OCR confidence: %.3f (threshold: %.3f)", avg_confidence, self._cfg.ocr_confidence_threshold)
            if avg_confidence < self._cfg.ocr_confidence_threshold:
                defects.append(Defect(
                    type=DEFECT_NO_CLEAR_TEXT,
                    description=(
                        f"Average OCR confidence {avg_confidence:.2f} is below threshold "
                        f"{self._cfg.ocr_confidence_threshold:.2f}. Text is not clearly readable."
                    ),
                ))

        return AnalysisResult(passed=not defects, defects=defects)
