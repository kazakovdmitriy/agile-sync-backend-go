# Makefile для Backend Go проекта

# Переменные
BINARY_NAME=agile_sync
MIGRATION_BINARY_NAME=migrate
MAIN_PATH=cmd/agile_sync/main.go
MIGRATION_PATH=cmd/migrate/main.go
SWAGGER_DOCS_PATH=docs
CONFIG_PATH=configs/config.yaml
MIGRATION_DIR=migrations
MOCKS_DIR=internal/mocks
SERVICE_INTERFACE_FILE=internal/service/interfaces.go
REPOSITORY_INTERFACE_FILE=internal/repository/interfaces.go

# Цвета для вывода
GREEN=\033[0;32m
RED=\033[0;31m
YELLOW=\033[0;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: all build run test clean deps swagger migrate help generate-mocks

all: deps swagger generate-mocks build

## Установка зависимостей
deps:
	@echo "$(YELLOW)Устанавливаем зависимости...$(NC)"
	go mod download
	go mod tidy

## Генерация Swagger документации
swagger:
	@echo "$(YELLOW)Генерируем Swagger документацию...$(NC)"
	@which swag > /dev/null || (echo "$(RED)swag не установлен. Устанавливаем...$(NC)" && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g $(MAIN_PATH) --output $(SWAGGER_DOCS_PATH) --parseDependency --parseInternal
	@echo "$(GREEN)Swagger документация сгенерирована в $(SWAGGER_DOCS_PATH)$(NC)"

## Генерация моков с помощью mockgen
generate-mocks:
	@echo "$(BLUE)Генерируем моки с помощью mockgen...$(NC)"
	@which mockgen > /dev/null || (echo "$(RED)mockgen не установлен. Устанавливаем...$(NC)" && go install go.uber.org/mock/mockgen@latest)

	@echo "$(BLUE)Создаем папку для моков...$(NC)"
	@mkdir -p $(MOCKS_DIR)

	@echo "$(BLUE)Генерируем моки для сервисов...$(NC)"
	@if [ -f "$(SERVICE_INTERFACE_FILE)" ]; then \
		mockgen -source=$(SERVICE_INTERFACE_FILE) -destination=$(MOCKS_DIR)/mock_service.go -package=mocks || echo "$(YELLOW)Предупреждение: Не удалось сгенерировать моки для сервисов$(NC)"; \
	else \
		echo "$(RED)Файл интерфейсов сервисов не найден: $(SERVICE_INTERFACE_FILE)$(NC)"; \
	fi

	@echo "$(BLUE)Генерируем моки для репозиториев...$(NC)"
	@if [ -f "$(REPOSITORY_INTERFACE_FILE)" ]; then \
		mockgen -source=$(REPOSITORY_INTERFACE_FILE) -destination=$(MOCKS_DIR)/mock_repository.go -package=mocks || echo "$(YELLOW)Предупреждение: Не удалось сгенерировать моки для репозиториев$(NC)"; \
	else \
		echo "$(RED)Файл интерфейсов репозиториев не найден: $(REPOSITORY_INTERFACE_FILE)$(NC)"; \
	fi

	@echo "$(GREEN)Генерация моков завершена$(NC)"

## Генерация моков для конкретного интерфейса
generate-mock:
	@echo "$(BLUE)Генерируем мок для конкретного интерфейса...$(NC)"
	@which mockgen > /dev/null || (echo "$(RED)mockgen не установлен. Устанавливаем...$(NC)" && go install go.uber.org/mock/mockgen@latest)
	@if [ -z "$(source)" ] || [ -z "$(interface)" ]; then \
		echo "$(RED)Использование: make generate-mock source=path/to/file.go interface=InterfaceName$(NC)"; \
		echo "$(RED)Пример: make generate-mock source=internal/service/interface.go interface=AuthService$(NC)"; \
		echo "$(RED)Пример: make generate-mock source=internal/repository/interface.go interface=UserRepository$(NC)"; \
		exit 1; \
	fi
	@mkdir -p $(MOCKS_DIR)
	@mockgen -source=$(source) -destination=$(MOCKS_DIR)/mock_$(shell basename $(source) .go)_$(interface).go -package=mocks $(interface)
	@echo "$(GREEN)Мок для интерфейса $(interface) из файла $(source) сгенерирован$(NC)"

## Генерация моков для всего пакета
generate-package-mocks:
	@echo "$(BLUE)Генерируем моки для всего пакета...$(NC)"
	@which mockgen > /dev/null || (echo "$(RED)mockgen не установлен. Устанавливаем...$(NC)" && go install go.uber.org/mock/mockgen@latest)
	@if [ -z "$(package)" ]; then \
		echo "$(RED)Использование: make generate-package-mocks package=path/to/package$(NC)"; \
		echo "$(RED)Пример: make generate-package-mocks package=internal/service$(NC)"; \
		exit 1; \
	fi
	@mkdir -p $(MOCKS_DIR)
	@mockgen -destination=$(MOCKS_DIR)/mock_$(shell basename $(package)).go -package=mocks $(package)
	@echo "$(GREEN)Моки для пакета $(package) сгенерированы$(NC)"

## Сборка основного приложения
build:
	@echo "$(YELLOW)Собираем приложение...$(NC)"
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Приложение собрано: bin/$(BINARY_NAME)$(NC)"

## Сборка инструмента миграций
build-migrate:
	@echo "$(YELLOW)Собираем инструмент миграций...$(NC)"
	go build -o bin/$(MIGRATION_BINARY_NAME) $(MIGRATION_PATH)
	@echo "$(GREEN)Инструмент миграций собран: bin/$(MIGRATION_BINARY_NAME)$(NC)"

## Запуск основного приложения
run: swagger
	@echo "$(YELLOW)Запускаем приложение...$(NC)"
	CONFIG_PATH=$(CONFIG_PATH) go run $(MAIN_PATH)

## Запуск в development режиме с автоперезагрузкой
dev: swagger generate-mocks
	@echo "$(YELLOW)Запускаем в development режиме с air...$(NC)"
	@which air > /dev/null || (echo "$(RED)air не установлен. Устанавливаем...$(NC)" && go install github.com/cosmtrek/air@latest)
	CONFIG_PATH=$(CONFIG_PATH) air

## Миграции - применение всех новых
migrate: build-migrate
	@echo "$(YELLOW)Применяем миграции...$(NC)"
	./bin/$(MIGRATION_BINARY_NAME) up

## Просмотр статуса миграций
migrate-status: build-migrate
	@echo "$(YELLOW)Проверяем статус миграций...$(NC)"
	./bin/$(MIGRATION_BINARY_NAME) status

## Откат последней миграции
migrate-down: build-migrate
	@echo "$(YELLOW)Откатываем последнюю миграцию...$(NC)"
	./bin/$(MIGRATION_BINARY_NAME) down

## Просмотр текущей версии миграций
migrate-version: build-migrate
	@echo "$(YELLOW)Проверяем текущую версию миграций...$(NC)"
	./bin/$(MIGRATION_BINARY_NAME) version

## Создание новой миграции (используя goose напрямую)
migrate-create:
	@echo "$(YELLOW)Создаем новую миграцию...$(NC)"
	@which goose > /dev/null || (echo "$(RED)goose не установлен. Устанавливаем...$(NC)" && go install github.com/pressly/goose/v3/cmd/goose@latest)
	@read -p "Введите название миграции: " name; \
	goose -dir $(MIGRATION_DIR) create $${name// /_} sql

## Сброс всех миграций (опасно - только для разработки)
migrate-reset: build-migrate
	@echo "$(RED)ВНИМАНИЕ: Это приведет к откату ВСЕХ миграций!$(NC)"
	@read -p "Вы уверены, что хотите продолжить? (y/N): " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo "$(YELLOW)Выполняем сброс всех миграций...$(NC)"; \
		./bin/$(MIGRATION_BINARY_NAME) reset; \
	else \
		echo "Операция отменена."; \
	fi

## Проверка целостности миграций
migrate-validate: build-migrate
	@echo "$(YELLOW)Проверяем целостность миграций...$(NC)"
	@which goose > /dev/null || (echo "$(RED)goose не установлен. Устанавливаем...$(NC)" && go install github.com/pressly/goose/v3/cmd/goose@latest)
	goose -dir $(MIGRATION_DIR) fix

## Очистка сгенерированных файлов
clean:
	@echo "$(YELLOW)Очищаем сгенерированных файлы...$(NC)"
	rm -rf bin/
	rm -rf $(SWAGGER_DOCS_PATH)/
	rm -rf $(MOCKS_DIR)/mock_*.go
	@echo "$(GREEN)Очистка завершена$(NC)"

## Запуск тестов (с моками)
test: generate-mocks
	@echo "$(YELLOW)Запускаем тесты...$(NC)"
	go test ./... -v

## Запуск тестов с покрытием
test-coverage: generate-mocks
	@echo "$(YELLOW)Запускаем тесты с покрытием...$(NC)"
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Отчет о покрытии: coverage.html$(NC)"

## Запуск тестов с перегенерацией моков
test-with-mocks: generate-mocks test

## Линтинг кода
lint:
	@echo "$(YELLOW)Проверяем код линтером...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(RED)golangci-lint не установлен. Устанавливаем...$(NC)" && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

## Форматирование кода
fmt:
	@echo "$(YELLOW)Форматируем код...$(NC)"
	go fmt ./...

## Валидация Swagger документации
validate-swagger: swagger
	@echo "$(YELLOW)Валидируем Swagger документацию...$(NC)"
	@which swagger > /dev/null || (echo "$(RED)swagger не установлен. Устанавливаем...$(NC)" && go install github.com/go-swagger/go-swagger/cmd/swagger@latest)
	swagger validate $(SWAGGER_DOCS_PATH)/swagger.yaml

## Полная сборка проекта (включая миграции и моки)
all-build: deps swagger generate-mocks build build-migrate
	@echo "$(GREEN)Полная сборка завершена!$(NC)"
	@echo "Основное приложение: bin/$(BINARY_NAME)"
	@echo "Инструмент миграций: bin/$(MIGRATION_BINARY_NAME)"

## Запуск миграций и приложения (для быстрого старта)
start: migrate run

## Просмотр помощи
help:
	@echo "$(GREEN)Доступные команды:$(NC)"
	@echo ""
	@echo "$(YELLOW)Основные:$(NC)"
	@echo "  make all                    - Установка зависимостей, генерация документации, моков и сборка"
	@echo "  make run                    - Запуск приложения"
	@echo "  make dev                    - Запуск в development режиме с автоперезагрузкой"
	@echo "  make build                  - Сборка основного приложения"
	@echo "  make start                  - Применить миграции и запустить приложение"
	@echo ""
	@echo "$(YELLOW)Документация:$(NC)"
	@echo "  make swagger                - Генерация Swagger документации"
	@echo "  make validate-swagger       - Валидация Swagger документации"
	@echo ""
	@echo "$(BLUE)Моки (gomock):$(NC)"
	@echo "  make generate-mocks         - Генерация всех моков"
	@echo "  make generate-mock          - Генерация мока для конкретного интерфейса"
	@echo "                               (использование: make generate-mock source=file.go interface=Name)"
	@echo "  make generate-package-mocks - Генерация моков для всего пакета"
	@echo "                               (использование: make generate-package-mocks package=path/to/package)"
	@echo ""
	@echo "$(YELLOW)Миграции:$(NC)"
	@echo "  make migrate                - Применить все новые миграции"
	@echo "  make migrate-status         - Просмотр статуса миграций"
	@echo "  make migrate-down           - Откат последней миграции"
	@echo "  make migrate-version        - Просмотр текущей версии миграций"
	@echo "  make migrate-create         - Создание новой миграции"
	@echo "  make migrate-reset          - Сброс ВСЕХ миграций (только для разработки!)"
	@echo "  make migrate-validate       - Проверка целостности миграций"
	@echo ""
	@echo "$(YELLOW)Тестирование:$(NC)"
	@echo "  make test                   - Запуск тестов (с автогенерацией моков)"
	@echo "  make test-with-mocks        - Перегенерировать моки и запустить тесты"
	@echo "  make test-coverage          - Запуск тестов с покрытием"
	@echo ""
	@echo "$(YELLOW)Разработка:$(NC)"
	@echo "  make lint                   - Проверка кода линтером"
	@echo "  make fmt                    - Форматирование кода"
	@echo ""
	@echo "$(YELLOW)Очистка:$(NC)"
	@echo "  make clean                  - Очистка сгенерированных файлов (бинарники, доки, моки)"
	@echo ""
	@echo "$(YELLOW)Другое:$(NC)"
	@echo "  make help                   - Показать эту справку"