package service

import (
	"WB_L0/app/repository/mocks"
	"WB_L0/domain"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetBackOrderHandlerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	api := &Apistruct{OrderRepo: mockRepo}

	expected := &domain.IncomingData{
		OrderUID:    "testo",
		TrackNumber: "ChevBistro",
		Entry:       "Prohod",
	}

	mockRepo.EXPECT().
		GetOrder("testo").
		Return(expected, nil).
		Times(1)

	// имитация запроса http
	r := httptest.NewRequest("GET", "/order?order_uid=testo", nil)
	w := httptest.NewRecorder()
	//  вызов проверяемой функции
	api.GetBackOrderhandler(w, r)

	//	проверка
	assert.Equal(t, http.StatusOK, w.Code)
	var response domain.IncomingData
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "testo", response.OrderUID)
}

func TestPut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	api := &Apistruct{OrderRepo: mockRepo}

	tested := &domain.IncomingData{
		OrderUID:    "testo",
		TrackNumber: "ChevBistro",
		Entry:       "Prohod",
	}
	jsData, err := json.Marshal(tested)

	assert.NoError(t, err)

	mockRepo.EXPECT().
		SaveOrder(gomock.Any()).
		Return(nil)

	r := httptest.NewRequest("POST", "/order", bytes.NewReader(jsData))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	api.PutData(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}
