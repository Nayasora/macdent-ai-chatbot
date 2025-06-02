ifeq ($(wildcard .env),)
	$(shell cp -n .env.example .env 2>/dev/null || true)
	ifneq ($(MAKECMDGOALS),)
		$(info .env файл не найден в корне проекта и он скопирован из .env.example)
		$(error Надо настроить .env файл и запустить make еще раз)
	endif
endif

include .env
export

MODE ?= development
COMPOSE_FILE=compose.$(MODE).yml
NETWORK_NAME=ai-chatbot

.PHONY: build up down restart logs clean help

build:
	docker compose -f $(COMPOSE_FILE) build

up:
	docker network inspect $(NETWORK_NAME) >/dev/null 2>&1 || docker network create $(NETWORK_NAME)
	docker compose -f $(COMPOSE_FILE) up

upd:
	docker compose -f $(COMPOSE_FILE) up -d

down:
	docker compose -f $(COMPOSE_FILE) down

stop:
	docker compose -f $(COMPOSE_FILE) stop

restart: down up

logs:
	docker compose -f $(COMPOSE_FILE) logs -f

rebuild: down build up

clean:
	docker compose -f $(COMPOSE_FILE) down -v
	docker system prune -f

sh:
	docker compose -f $(COMPOSE_FILE) exec client sh

shr:
	docker compose -f $(COMPOSE_FILE) exec -u root client sh

cmg:
	docker rmi -f macdent-ai-chatbot-nginx macdent-ai-client-development