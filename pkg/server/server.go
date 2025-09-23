package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	// "github.com/dokedza/WB_L0/domain"
	api "github.com/dokedza/WB_L0/app/handler"
	"github.com/dokedza/WB_L0/app/service"
	kf "github.com/dokedza/WB_L0/pkg/kafkain"
)

// функция создания сервера
func Run(c *service.Apistruct) error {
	port := 7540

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
	fmt.Printf("Сервер запущен на порту: http://localhost:%d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func checkConsume() {
	for {
		kf.ConsumeMessage()
	}
}
