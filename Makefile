.PHONY: run import

ifneq (,$(wildcard .env))
    include .env
    export
endif

run: # Runs the API server.
	go run cmd/api-server/main.go

import: # Imports movies from a CSV file into the database. eg: make import ARGS="-file=movies.csv -size=200"
	go run cmd/import-movies/main.go $(ARGS)
