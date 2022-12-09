.POSIX:
.SUFFIXES:

.PHONY: default
default: dist

venv:
	python3 -m venv venv
	./venv/bin/python3 -m pip install -r requirements-dev.txt

.PHONY: build
build: clean
	python3 -m pip install -q -r requirements.txt --only-binary=:all: --platform any --target build/
	rm -fr build/*.dist-info/
	cp manage.py build/
	cp -r project build/
	cp -r digimontcg build/
	shiv --compressed --site-packages build -p "/usr/bin/env python3" -o digimontcg.pyz -e manage:main

.PHONY: static
static: venv
	./venv/bin/python3 manage.py collectstatic --no-input

.PHONY: dist
dist: build static
	rm -fr dist/
	mkdir dist/
	cp *.pyz dist/
	cp -r static dist/

.PHONY: package
package: dist
	nfpm package -p deb -t digimontcg.deb

.PHONY: test
test: venv
	./venv/bin/python3 manage.py test

.PHONY: format
format: venv
	./venv/bin/black manage.py digimontcg/ project/

.PHONY: clean
clean:
	rm -fr *.deb *.pyz build/ dist/ static/
