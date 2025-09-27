package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	db "WB_L0/app/repository"

	order "WB_L0/domain"
)

type DataResp struct {
	Data []*order.IncomingData `json:"data"`
}
type Apistruct struct {
	Bdstruct *db.Bdstruct
}

// возврат данных по UID
func (ap *Apistruct) GetBackOrderhandler(w http.ResponseWriter, r *http.Request) {

	Uid := r.URL.Query().Get("order_uid")
	if Uid == "" {
		jsonWriter(w, map[string]string{"error": "отсутствует UID"})
		return
	}
	DataResp, err := ap.Bdstruct.GiveBackOrdrerData(Uid)
	if err != nil {
		jsonWriter(w, map[string]any{"error": err})
		return
	}
	jsonData, err := json.Marshal(DataResp)
	if err != nil {
		jsonWriter(w, map[string]any{"error": fmt.Errorf("Ошибка отправки: %w", err).Error()})
		return
	}
	w.Write(jsonData)
}

// доставление данных в бд
func (ap *Apistruct) PutData(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var dataResp order.IncomingData

	body, err := io.ReadAll(r.Body)

	if err != nil {
		jsonWriter(w, map[string]any{"error": fmt.Errorf("Ошибка входящих данных: %w", err).Error()})
		return
	}
	err = json.Unmarshal(body, &dataResp)
	if err != nil {
		jsonWriter(w, map[string]string{"error": fmt.Errorf("Ошибка декодирования: %w", err).Error()})
		return
	}
	err = db.DataHasArrivedInOrders(&dataResp, ap.Bdstruct)
	if err != nil {
		jsonWriter(w, map[string]string{"error": fmt.Errorf("Ошибка связи с таблицами: %w", err).Error()})
		return
	}
}

// функция созданная исключительно для удобства написания ответов
func jsonWriter(w http.ResponseWriter, data any) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}
