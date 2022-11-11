.POSIX:
.SUFFIXES:

.PHONY: default
default: dist

venv:
	python3 -m venv venv
	./venv/bin/python3 -m pip install -r requirements-dev.txt

.PHONY: test
test: venv
	./venv/bin/pytest

.PHONY: format
format:
	./venv/bin/black manage.py digimontcg/ project/
