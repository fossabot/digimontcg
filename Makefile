.POSIX:
.SUFFIXES:

.PHONY: default
default: dist

venv:
	python3 -m venv venv
	./venv/bin/python3 -m pip install -r requirements-dev.txt

.PHONY: build
build: 
	rm -fr build/
	python3 -m pip install -q -r requirements.txt --target build/
	rm -fr build/*.dist-info/
	cp -r main.py build/__main__.py
	python3 -m zipapp -c -p "/usr/bin/env python3" -o dmtcgo build/

.PHONY: dist
dist: build
	rm -fr dist/
	mkdir dist/
	cp dmtcgo dist/
	cp -r data dist/

.PHONY: test
test: venv
	./venv/bin/pytest

.PHONY: format
format:
	./venv/bin/black main.py scrape.py

.PHONY: clean
clean:
	rm -fr dmtcgo dist/ build/ venv/
