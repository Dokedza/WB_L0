package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	api "WB_L0/app/handler"
	"WB_L0/app/service"
	kf "WB_L0/pkg/kafkain"
)

// функция создания сервера
func Run(c *service.Apistruct) error {
	port := 8080

	envPort := os.Getenv("TODO_PORT")
	if envPort != "" {
		p, err := strconv.Atoi(envPort)
		if err == nil {
			port = p
		}
	}
	//постоянная проверка на наличие новых сообщений
	go checkConsume()

	api.Init(c)
	http.Handle("/", http.FileServer(http.Dir("web")))
	log.Printf("Сервер запущен на порту: http://localhost:%d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func checkConsume() {
	for {
		kf.ConsumeMessage()
	}
}
