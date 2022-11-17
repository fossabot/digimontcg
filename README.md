# digimontcg
They are the champions

## Setup
This project depends on the [Python programming language](https://www.python.org/) and uses the [Django web framework](https://www.djangoproject.com/).
I like to use a [POSIX-compatible Makefile](https://pubs.opengroup.org/onlinepubs/9699919799.2018edition/utilities/make.html) to facilitate the various project operations but traditional commands will work just as well.

## Install
If you are unfamiliar with [virtual environments](https://docs.python.org/3/library/venv.html), I suggest taking a brief moment to learn about them and set one up.
The Python docs provide a great [tutorial](https://docs.python.org/3/tutorial/venv.html) for getting started with virtual environments and packages.

This project's dependencies can be installed via pip:
```
pip install -r requirements-dev.txt
```

## Running
To start the web server:
```
python3 manage.py runserver
```
