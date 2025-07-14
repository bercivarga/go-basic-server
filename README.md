# Go Basic Server

This is a basic Go server that provides a RESTful API for user authentication and management.
It uses Goose for managing migrations, SQLC for generating SQL queries, JWTs for authentication, and Air for local development.

This project is meant to be a starting point for building backends. Adding additional architecture features like Docker, caching, load balancing, secret management etc. should be left to the individual projects.

## Setup
This project has a Makefile with some started commands that can get a project started:
```shell
make init # Initialize project after cloning
make run # Run with air (live reload)

make test # Run tests
make vet # Run vet for Go and SQLC

make generate # Generate SQL queries with SQLC

make up # Run Goose up migrations
make down # Run Goose down migrations
make create # Create a new migration file
```

Note: you will need Go v1.24.2 or higher.
