IMAGE_NAME = tg-proxy
COMPOSE_FILE = docker-compose.yml

.PHONY: dbuild
dbuild: 
	docker build -t ${IMAGE_NAME} .

.PHONY: compose_run
compose_run: 
	docker compose --env-file ./.env up --remove-orphans

.PHONY: dbuildrun
dbuildrun: dbuild compose_run

.PHONY: run
run: 
	go run main.go