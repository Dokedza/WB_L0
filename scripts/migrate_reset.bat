@echo off
echo ПЕРЕСОЗДАНИЕ БАЗЫ ДАННЫХ!
echo ВСЕ ДАННЫЕ БУДУТ УДАЛЕНЫ!
pause

echo Удаляем базу данных...
psql postgres -c "drop database if exists wb_orders"

echo Создаем новую базу...
createdb wb_orders

echo Применяем миграции...
goose -allow-missing -dir migrations postgres "user=postgres password=5037 dbname=wb_orders sslmode=disable" up

echo База данных пересоздана!
pause