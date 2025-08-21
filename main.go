package main

import (
	"fmt"
	database "go1f/pkg/database"
	server "go1f/pkg/server"

	"github.com/dokedza/WB_L0/pkg/api"
	//"github.com/dokedza/WB_L0/pkg/kafkain"
)

func main() {
	//kafkain.ConsumeMessage()

	//инициация всех таблиц
	err := database.Init()
	if err != nil {
		fmt.Println(err)
	}

	myCache := database.New()
	myb := &database.Bdstruct{Cache: myCache}
	c := &api.Apistruct{Bdstruct: myb}
	myCache.CacheInit()

	//запуск сервера, обработка ошибок
	err = server.Run(c)
	if err != nil {
		panic(err)
	}
}
