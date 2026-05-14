.PHONY: install run update-protos

install:
	pip install -e .

run:
	python -m document_analyzer.main

# Re-download vendored proto stubs from the canonical proto repo tag
update-protos:
	python scripts/update_protos.py
