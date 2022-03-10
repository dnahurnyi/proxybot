include linting.mk

IMAGE_NAME = tg-proxy
COMPOSE_FILE = docker/docker-compose.yml

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