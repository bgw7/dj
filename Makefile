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


psql: ## Start a psql session with local DB
	psql

db-start: db-stop ## Start postgres service and create database
	make -C ./db start

db-stop: ## Stop postgres service | shut down the database
	-make -C ./db stop

api-start: ## Start the ./api application
	make -C ./api start

ui-start: ui-stop ## Start the ./ui nginx server
	make -C ./ui start

ui-stop: ## Stop the ./ui service
	make -C ./ui stop

start: db-start api-start ## Start DB, API
	print db, api ready!
	
stop: db-stop ui-stop
