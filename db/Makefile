#!/usr/bin/make -f

DB_SUPER_USER=${USER}

ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.DEFAULT_GOAL := help

.PHONY: help

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\n\nWhere <target> is one of:\n"} /^[$$()% a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

sql-init: ## Run *.sql files in ./schemas directory
	./init.sh

psql: ## Start terminal session with postgres service using .env Vars & psql CLI
	psql

sql-adhoc: ## Run SQL statements in ./testing/sql/adhoc.sql
	psql -U $(DB_SUPER_USER) -f ./testing/sql/adhoc.sql

stop: ## Stop postgres service | shut down the database
	-pg stop

create-db: stop  ## Start postgres service and create database
	-pg start
	-dropdb --username=$(DB_SUPER_USER) $(PGDATABASE)
	-createdb --username=$(DB_SUPER_USER) $(PGDATABASE)

start: stop create-db sql-init ## Starts postgres service, creates database, runs sql-init
	@echo local postgres running with $(PGDATABASE) database

## %: this target is needed to allow args to be passed to tagets that are using: $(filter-out $@, $(MAKECMDGOALS))
%:
	@true
