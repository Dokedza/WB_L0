package main

import (
	"WB_L0/app/migrator"
	server "WB_L0/pkg/server"
	"database/sql"
	"log"

	database "WB_L0/app/repository"

	api "WB_L0/app/service"
	ch "WB_L0/pkg/cache"
)

func main() {
	// подключение к бд
	ConnStr := "host=postgres password=5037 user=postgres dbname=WB_L0 sslmode=disable"
	db, err := sql.Open("postgres", ConnStr)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()
	// автомиграции
	mgr := migrator.New(db)
	err = mgr.RunMigrations()
	if err != nil {
		log.Fatal("Ошибка миграций", err)
	}
	//нинициализация кэша
	myCache := ch.New()
	myb := &database.Bdstruct{Cache: myCache}
	c := &api.Apistruct{Bdstruct: myb}

	c.Bdstruct.Cache.CacheInit()

	// myCache.CacheInit()
	// запуск сервера, обработка ошибок
	err = server.Run(c)
	if err != nil {
		panic(err)
	}
}
