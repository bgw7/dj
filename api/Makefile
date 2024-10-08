#!/usr/bin/make -f

DB_SUPER_USER=${USER}

ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.DEFAULT_GOAL := help

.PHONY: help start-db

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\n\nWhere <target> is one of:\n"} /^[$$()% a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

fmt: # Run go fmt
	go fmt ./...

test-unit: fmt ## Run unit test
	-rm ./coverage.out
	-go test -v ./... -coverprofile=coverage.out -race
	go tool cover -html=coverage.out

test-acceptance:
	-go test ./testing/acceptance_test.go -test.v -test.run ^TestFeatures --godog.tags="@essential"

start: fmt  ## Run this Go-lang application
	go run .

db-start: ## Start the relational datastore
	-make -C ../db start

post: ## create reservation via curl
	curl -v localhost:9999/api/reservations -d "@./testing/testdata/input/createReservation.json"
