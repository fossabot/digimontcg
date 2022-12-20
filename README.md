# digimontcg
They are the champions

## Setup
This project depends on the [Python programming language](https://www.python.org/) and uses the [Django web framework](https://www.djangoproject.com/).
I like to use a [POSIX-compatible Makefile](https://pubs.opengroup.org/onlinepubs/9699919799.2018edition/utilities/make.html) to facilitate the various project operations but traditional commands will work just as well.

## Local Development
### Dependencies
If you are unfamiliar with [virtual environments](https://docs.python.org/3/library/venv.html), I suggest taking a brief moment to learn about them and set one up.
The Python docs provide a great [tutorial](https://docs.python.org/3/tutorial/venv.html) for getting started with virtual environments and packages.

This project's dependencies can be installed via pip:
```
pip install -r requirements-dev.txt
```

### Services
This project uses [PostgreSQL](https://www.postgresql.org/) for persistent storage and [Redis](https://redis.io/) for queuing background tasks and general caching.
To develop locally, you'll need to run these services locally somehow or another.
I find [Docker](https://www.docker.com/) to be a nice tool for this but you can do whatever works best.

The following command starts the necessary containers:
```
docker compose up -d
```

These containers can be stopped via:
```
docker compose down
```

### Running
To start the web server:
```
python3 manage.py runserver
```

## Design
The design for this project is heavily inspired by the [Pokémon TCG Developers](https://pokemontcg.io/) website.
In my opinion, the work they have done sets a standard of quality that all open-source API projects should work toward.

## Legal
This project is not produced, endorsed, supported, or affiliated with Bandai Co., Ltd.
