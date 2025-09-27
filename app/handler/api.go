package api

import (
	"net/http"

	ap "WB_L0/app/service"
)

func Init(c *ap.Apistruct) {
	http.HandleFunc("/api/task/getBackOrderhandler", corsMiddleware(c.GetBackOrderhandler))
	http.HandleFunc("/api/task/PutData", c.PutData)
}
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                   // Разрешает все источники
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Разрешённые методы
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // Разрешённые заголовки

		// Обработка preflight-запросов
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
