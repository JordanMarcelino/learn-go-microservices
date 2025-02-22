.PHONY: all build compose-up

all: build compose-up

dirs := ./auth-service ./product-service ./order-service ./mail-service ./gateway

migrate_dirs := ./auth-service ./product-service ./order-service

build:
	@echo "Building all services..."
	@for dir in $(dirs); do \
		echo "Building $$dir..."; \
		cd $$dir && make build && cd -; \
	done

migrateup:
	@echo "Migrating all services..."
	@for dir in $(migrate_dirs); do \
		echo "Migrating $$dir..."; \
		DB_HOST=localhost; \
		if [ "$$dir" = "./auth-service" ]; then \
			DB_PORT=5000; \
		elif [ "$$dir" = "./product-service" ]; then \
			DB_PORT=5001; \
		elif [ "$$dir" = "./order-service" ]; then \
			DB_PORT=5002; \
		fi; \
		migrate -path $$dir/db/migrations/ -database "postgresql://postgres:postgres@$$DB_HOST:$$DB_PORT/api_db?sslmode=disable" -verbose up; \
	done

compose-up:
	@docker network create --driver bridge production
	@docker compose -f deployments/compose-pgpool.yaml \
		-f deployments/compose-redis.yaml \
		-f deployments/compose-kafka.yaml \
		-f deployments/compose-rabbitmq.yaml \
		up -d --build
	@make migrateup
	@docker compose -f deployments/compose-api.yaml up -d --build

compose-down:
	@docker compose -f deployments/compose-pgpool.yaml \
		-f deployments/compose-redis.yaml \
		-f deployments/compose-kafka.yaml \
		-f deployments/compose-rabbitmq.yaml \
		down -v
	@docker compose -f deployments/compose-api.yaml down -v
	@docker network rm production
