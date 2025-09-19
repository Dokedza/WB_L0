package main

import (
	"fmt"

	database "go1f/pkg/database"
	server "go1f/pkg/server"

	"github.com/dokedza/WB_L0/pkg/api"
	ch "github.com/dokedza/WB_L0/pkg/cache"
)

func main() {

	//инициация всех таблиц
	err := database.Init()
	if err != nil {
		fmt.Println(err)
	}

	//нинициализация кэша
	myCache := ch.New()
	myb := &database.Bdstruct{Cache: myCache}
	c := &api.Apistruct{Bdstruct: myb}

	// c.Bdstruct.Cache.CacheInit()

	myCache.CacheInit()
	// запуск метода перезаписи кэша раз в 1 час
	myCache.StartAuthoRefresh()
	// мягкая остановка перезаписи кэша при остановке программы
	defer myCache.StopAuthoRefresh()
	// запуск сервера, обработка ошибок
	err = server.Run(c)
	if err != nil {
		panic(err)
	}
}
