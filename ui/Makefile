ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\n\nWhere <target> is one of:\n"} /^[$$()% a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


start: stop  ## start nginx service, pid inside nginx.pid
	nginx -p ./nginx/conf.d -c default.conf
	> ./nginx/conf.d/logs/access.log
	@echo running on port 7770
	tail -f ./nginx/conf.d/logs/access.log

stop: ## kill nginx pid
	-kill $(shell cat ./nginx/conf.d/logs/nginx.pid)
