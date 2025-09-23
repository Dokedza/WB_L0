@echo off
echo Применяем миграции...
goose -dir migrations postgres "user=postgres password=5037 dbname=wb_orders sslmode=disable" up
echo Готово!
pause