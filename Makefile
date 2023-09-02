# User defined variables
USER := $(whoami)
REDIS_VOLUME := ~/OrbStack/docker/volumes/symlet-backend_redis-data/
# Go variables
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
BINARY_NAME := myapp

# Docker variables
DOCKER_COMPOSE := docker-compose
DOCKER_COMPOSE_FILE := docker-compose.yaml

.PHONY: all build test clean up down clean-redis

all: test build up

build:
	@$(GOBUILD) -o $(BINARY_NAME) -v ./...

test:
	@$(GOTEST) -v ./...

clean:
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)

up:
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d

down:
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down
clean-redis:
	@rm -f $(REDIS_VOLUME)