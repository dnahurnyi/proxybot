include linting.mk

IMAGE_NAME = tg-proxy
COMPOSE_FILE = docker/docker-compose.yml
kill_services = (echo "External postgres - nothing to kill")
remove_not_running = (echo "External postgres - nothing to remove")
run_docker = (echo "External postgres - nothing to run")
integration_tests_tags = "integration"
ifeq ($(origin DB_URL), undefined)
	POSTGRES_HOST = localhost
	POSTGRES_PORT = 5432
	POSTGRES_USER = user
	POSTGRES_PASSWORD = password
	POSTGRES_DB = proxybot-db
	DB_URL = postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable
	kill_services = (docker kill postgres 2> /dev/null || true)
	remove_not_running = (docker rm $$(docker ps -a -q) || true)
	run_docker = (docker run -d --name postgres -p $(POSTGRES_PORT):$(POSTGRES_PORT) \
	-e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_DB=$(POSTGRES_DB) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) postgres:11 && sleep $(POSTGRES_WAIT_TIME))
endif

.PHONY: dbuild
dbuild: 
	docker build -t ${IMAGE_NAME} .

.PHONY: compose_run
compose_run: 
	docker compose -f docker/docker-compose.yml up --remove-orphans

.PHONY: dbuildrun
dbuildrun: dbuild compose_run

.PHONY: run
run: 
	go run main.go

.PHONY: mocks
mocks:
	mockgen -destination=bot/mock/repository.go -package=mock github.com/dnahurnyi/proxybot/bot Repository
	mockgen -destination=bot/mock/client.go -package=mock github.com/dnahurnyi/proxybot/bot Client
	mockgen -destination=bot/mock/id_gen.go -package=mock github.com/dnahurnyi/proxybot/bot IDGenerator

.PHONY: unit_test
unit_test:
	go test -mod=vendor --tags=unit -count=1 -v -cover  -timeout 300s ./...

.PHONY: ci
ci: lint unit_test integration_test

.PHONY: integration_test
integration_test:
	$(call kill_services)
	$(call remove_not_running)
	$(call run_docker)
	export DB_URL="$(DB_URL)" && \
	go test -mod=vendor --tags=$(integration_tests_tags) -count=1 -v -cover ./storage/...
	$(call kill_services)
	$(call remove_not_running)