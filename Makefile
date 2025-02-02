PROJECT_NAME := himaya-go
user := $(shell id -u)
group := $(shell id -g)

go_version := $(shell cat go.mod | grep '^go' | sed -E "s/go\s//")

PORT?=3000
dc = USER_ID=$(user) GROUP_ID=$(group) GO_VERSION="$(go_version)" PORT=$(PORT) COMPOSE_PROJECT_NAME=$(PROJECT_NAME) docker-compose -p $(PROJECT_NAME) -f docker-compose.yaml $(1)$(2)
dr = @$(call dc, run --rm $(1)$(2))
drt = @$(call dc, run --rm --env APP_ENV=test $(1)$(2))
COMPOSE_OPTIONS ?=
de = @$(call dc, exec $(COMPOSE_OPTIONS) $(1)$(2))

args = $(shell arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}})

entry := cmd/server/main.go
include .env

.DEFAULT_GOAL := help

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

TARGET_MAX_CHAR_NUM=25
## Show this help message
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			sub(/:/, "", helpCommand); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.PHONY: dc
## Run docker-compose
dc:
	@$(call dc, $(args))

.PHONY: pre-dev
pre-dev: .env
	@$(call dc, up -d database redis)

dev_services = backend node meilisearch
.PHONY: dev
## Start development containers in foreground
dev: pre-dev history/go/fish history/node/fish
	EXTERNAL_PORT=$(PORT) $(call dc, up $(if $(args),$(args),$(dev_services)))

.PHONY: up
## Start containers in background
up:
	EXTERNAL_PORT=$(PORT) $(call dc, up -d $(args))

.PHONY: stop
## Stop containers
stop:
	EXTERNAL_PORT=$(PORT) $(call dc, stop $(args))

.PHONY: down
## Stop and remove all containers (including runners)
down:
	$(call dc, down --remove-orphans)

.PHONY: restart
## Recreate container
restart:
	$(call dc, up --force-recreate -d $(args))

### Building and installing ###

.PHONY: setup
## Setup the project (build images, install deps, etc)
setup: down .env _bundless history/go/fish history/node/fish
	make docker-build
	$(call dr, node yarn install)
	$(call dr, node yarn run build)
	$(call dr, backend go mod tidy)
	-make import-base-data
	@echo "✅ ${GREEN}Build complete!${RESET}"

.env: .env.example
	@test -f $@ || (echo -e "APP_KEY=$(shell openssl rand -base64 32)"; cat $@.example) > $@

.PHONY: history/%
history/%:
	@if [ -d docker/$@ ]; then rm -rfv docker/$@; fi
	@if [ ! -f docker/$@ ]; then touch docker/$@; fi

.PHONY: _bundless
_bundless:
	-docker volume rm $(PROJECT_NAME)_go_pkg

.PHONY: docker-build
## Build docker images
docker-build:
	$(call dc, build --pull $(args))

.PHONY: build
## Build production image
build:
	$(call dc, build production)

.PHONY: import-base-data
## Import data from a dump files
import-base-data:
	@for file in dumps/*.sql.gz; do \
		file_name=$$(basename $$file .sql.gz); \
		table_name=$$(echo $$file_name | sed -E 's/^[0-9]+-//'); \
		$(dc) exec -T database psql -U $(DB_USERNAME) -d $(DB_NAME) -c "TRUNCATE $$table_name CASCADE"; \
		gunzip -c $$file | $(dc) exec -T database psql \
						-U $(DB_USERNAME) -d $(DB_NAME) \
						-c "COPY $$table_name FROM STDIN"; \
	done

.PHONY: fish/%
## Get some fish (backend, node)
fish/%:
	$(dr) $* /usr/bin/fish

.PHONY: docs
## Generate swagger docs
docs:
	$(call de, backend swag fmt)
	$(call de, backend swag init -g $(entry) --parseDependency --parseInternal -q)
	$(call de, node swagger2openapi docs/swagger.yaml -o docs/openapi.yaml)
	$(call de, node sh scripts/fix_openapi.sh docs/openapi.yaml)
	$(call de, node yarn --silent o2ts docs/openapi.yaml)
	@mv -v docs/openapi.ts static/api.ts


.PHONY: psql
## DUH!
psql:
	$(call de, database psql -U $(DB_USERNAME) -d $(DB_NAME))

.PHONY: test
## Run test using ginkgo
test:
	$(call drt, backend ginkgo -r)

.PHONY: guard
## Watch and run tests
guard:
	$(call drt, backend ginkgo watch -r)

.PHONY: touch
## Touch entrypoint to restart server
touch:
	touch $(entry)
