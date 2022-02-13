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