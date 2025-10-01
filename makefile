.PHONY: generate-mocks test clean

# Генерация всех моков
generate-mocks:
	@echo "Generating mocks..."
	mockgen -source=app/repository/repository.go -destination=app/repository/mocks/repository_mock.go -package=mocks
	mockgen -source=app/service/service.go -destination=app/service/mocks/service_mock.go -package=mocks
	mockgen -source=pkg/cache/cache.go -destination=pkg/cache/mocks/cache_mock.go -package=mocks
	@echo "Mocks generated successfully!"

# Генерация мока для конкретного интерфейса
mock-repository:
	mockgen -source=app/repository/repository.go -destination=app/repository/mocks/repository_mock.go -package=mocks

mock-service:
	mockgen -source=app/service/service.go -destination=app/service/mocks/service_mock.go -package=mocks

mock-cache:
	mockgen -source=pkg/cache/cache.go -destination=pkg/cache/mocks/cache_mock.go -package=mocks

# Запуск тестов с моками
test: generate-mocks
	go test ./... -v

# Очистка сгенерированных моков
clean-mocks:
	rm -f app/repository/mocks/*.go
	rm -f app/service/mocks/*.go
	rm -f pkg/cache/mocks/*.go

# Установка mockgen (если не установлен)
install-mockgen:
	go install github.com/golang/mock/mockgen@latest

# Помощь
help:
	@echo "Available targets:"
	@echo "  generate-mocks - Generate all mocks"
	@echo "  mock-repository - Generate repository mock only"
	@echo "  mock-service - Generate service mock only"
	@echo "  mock-cache - Generate cache mock only"
	@echo "  test - Run tests with mocks"
	@echo "  clean-mocks - Remove generated mocks"
	@echo "  install-mockgen - Install mockgen tool"