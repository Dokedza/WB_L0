package repository_test

import (
	"database/sql"
	"testing"

	"WB_L0/app/repository/mocks"
	"WB_L0/domain"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetSuccess(t *testing.T) {
	// контроллер мока
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// моковая фальвка
	mockRepo := mocks.NewMockOrderRepository(ctrl)

	// то что ожидается вернуть
	expectedOrder := &domain.IncomingData{
		OrderUID:    "testo",
		TrackNumber: "ChevBistro",
		Entry:       "Prohod",
	}
	// ожидаемые значения
	mockRepo.EXPECT().
		GetOrder("testo").
		Return(expectedOrder, nil).
		Times(1)
	// вызов для начала проверки
	result, err := mockRepo.GetOrder("testo")
	// проверка
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testo", result.OrderUID)

}

func TestGetNotFound(t *testing.T) {
	// контроллер мока
	// всё как и в первом тесте
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// моковая фальвка
	mockRepo := mocks.NewMockOrderRepository(ctrl)

	// ожидания
	mockRepo.EXPECT().
		GetOrder("ну такого уид точно не может быть").
		Return(nil, sql.ErrNoRows).
		Times(1)

	// выполнение
	result, err := mockRepo.GetOrder("ну такого уид точно не может быть")

	// проверки
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sql.ErrNoRows, err) //проверка на конкретную ошибку
}

func TestSaveSuccess(t *testing.T) {
	// контроллер мока
	// всё как и в первом тесте
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// моковая фальвка
	mockRepo := mocks.NewMockOrderRepository(ctrl)

	ordertest := &domain.IncomingData{
		OrderUID:    "testo",
		TrackNumber: "ChevBistro",
		Entry:       "Prohod",
	}

	// ожидания
	mockRepo.EXPECT().
		SaveOrder(ordertest).
		Return(nil).
		Times(1)

	// выполнение
	err := mockRepo.SaveOrder(ordertest)
	// проверка
	assert.NoError(t, err)
}

func TestSaveErr(t *testing.T) {
	// контроллер мока
	// всё как и в первом тесте
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// моковая фальвка
	mockRepo := mocks.NewMockOrderRepository(ctrl)
	orderTest := &domain.IncomingData{
		OrderUID: "test",
	}
	// ожидания
	mockRepo.EXPECT().
		SaveOrder(orderTest).
		Return(sql.ErrConnDone).
		Times(1)
	err := mockRepo.SaveOrder(orderTest)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
}

// func TestNilOrderUID(t *testing.T){
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	// моковая фальвка
// 	mockRepo := mocks.NewMockOrderRepository(ctrl)

// 	ordertestNil := &domain.IncomingData{
// 		OrderUID:    "",
// 		TrackNumber: "ChevBistro",
// 		Entry:       "Prohod",
// 	}

// 	mockRepo.EXPECT().
// 	SaveOrder(ordertestNil).
// 	Return(sql.ErrNoRows).
// 	Times(1)

// 	err:= mockRepo.SaveOrder(ordertestNil)

// 	assert.Error(t,err)
// }
