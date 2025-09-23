# Переменные
LOCAL_DB_NAME = wb_orders
LOCAL_DB_DSN = "user=postgres password=5037 dbname=$(LOCAL_DB_NAME) sslmode=disable"

# Пересоздание БД с миграциями
migrate-reset:
	psql postgres -c "drop database if exists $(LOCAL_DB_NAME)"
	createdb $(LOCAL_DB_NAME)
	goose -allow-missing -dir migrations postgres $(LOCAL_DB_DSN) up
	@echo "База данных пересоздана и миграции применены!"

# Только применение миграций
migrate-up:
	goose -dir migrations postgres $(LOCAL_DB_DSN) up
	@echo "Миграции применены!"

# Откат последней миграции
migrate-down:
	goose -dir migrations postgres $(LOCAL_DB_DSN) down
	@echo "Последняя миграция откачена!"

# Статус миграций
migrate-status:
	goose -dir migrations postgres $(LOCAL_DB_DSN) status

# Создание новой миграции
migrate-create:
	goose -dir migrations create $(name) sql
	@echo "Создана новая миграция: $(name)"

test-migrations:
    # Запуск БД
    docker-compose -f deployment/docker-compose.yml up -d postgres
    sleep 5
    # Применение миграций
    goose -dir migrations postgres $(LOCAL_DB_DSN) up
    # Проверка
	goose -dir migrations postgres $(LOCAL_DB_DSN) status
    # Откат
    goose -dir migrations postgres $(LOCAL_DB_DSN) down
    goose -dir migrations postgres $(LOCAL_DB_DSN) up
    echo "Миграции протестированы успешно!"