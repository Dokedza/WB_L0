package main

import (
	server "WB_L0/pkg/server"

	database "WB_L0/app/repository"

	api "WB_L0/app/service"
	ch "WB_L0/pkg/cache"
)

func main() {

	//нинициализация кэша
	myCache := ch.New()
	myb := &database.Bdstruct{Cache: myCache}
	c := &api.Apistruct{Bdstruct: myb}

	// c.Bdstruct.Cache.CacheInit()

	myCache.CacheInit()
	// запуск метода перезаписи кэша раз в 1 час
	// myCache.StartAuthoRefresh()
	// мягкая остановка перезаписи кэша при остановке программы
	// defer myCache.StopAuthoRefresh()
	// запуск сервера, обработка ошибок
	err := server.Run(c)
	if err != nil {
		panic(err)
	}
}
