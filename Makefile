ENV_LOCAL=./envs/local.env
# Загружаем переменные из .env-файла
include $(ENV_LOCAL)
export $(shell sed 's/=.*//' $(ENV_LOCAL))

# Переменные
MIGRATE = migrate
DB_URL = postgres://postgres:root@localhost:5432/postgres?sslmode=disable
MIGRATIONS_DIR = ./migrations

# Создание новой миграции
create-migration:
	@read -p "Введите название миграции: " name; \
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $${name}

# Применение миграций
migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database $(DB_URL) up

# Docker команды
docker-up:
	docker-compose --file $(COMPOSE_FILE) up -d

docker-down:
	docker-compose --file $(COMPOSE_FILE) down

docker-logs:
	docker-compose --file $(COMPOSE_FILE) logs -f

docker-restart:
	docker-compose --file $(COMPOSE_FILE) down && docker-compose --file $(COMPOSE_FILE) up -d

