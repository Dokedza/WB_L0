@echo off
chcp 65001 > nul
echo ========================================
echo   ТЕСТИРОВАНИЕ МИГРАЦИЙ
echo ========================================

echo [1/6] Проверяем наличие docker-compose.yml...
if not exist "deployment\docker-compose.yml" (
    echo ОШИБКА: Файл deployment\docker-compose.yml не найден!
    pause
    exit /b 1
)

echo [2/6] Запускаем PostgreSQL...
docker-compose -f deployment\docker-compose.yml up -d postgres

echo [3/6] Ждем 15 секунд пока БД запустится...
timeout /t 15

echo [4/6] Создаем базу данных если не существует...
psql -U postgres -c "CREATE DATABASE WB_L0;" 2>nul || echo База уже существует

echo [5/6] Применяем миграции...
goose -dir migrations postgres "user=postgres password=5037 dbname=WB_L0 sslmode=disable" up
if %errorlevel% neq 0 (
    echo ОШИБКА: Не удалось применить миграции!
    echo Проверь: 
    echo - Запущен ли Docker?
    echo - Запущен ли PostgreSQL в Docker?
    pause
    exit /b 1
)

echo [6/6] Проверяем статус...
goose -dir migrations postgres "user=postgres password=5037 dbname=WB_L0 sslmode=disable" status

echo   ТЕСТИРОВАНИЕ ЗАВЕРШЕНО УСПЕШНО!
pause