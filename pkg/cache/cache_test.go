package cache

import (
	"testing"

	"WB_L0/domain"

	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	// тест для получения

	cache := New()
	require.NotNil(t, cache, "Кэш должен быть создан")
	testOrder := domain.IncomingData{
		OrderUID:    "OneRing123",
		TrackNumber: "Mordor_Track",
		Entry:       "Rivendell",
		Locale:      "elvish",
		Delivery: domain.Delivery{
			Name:    "Frodo",
			Phone:   "+1234567788",
			Zip:     "30189",
			City:    "Shir",
			Address: "Bag end, Hobbiton",
			Region:  "Ennorat",
			Email:   "frodo123.gmail.mdl",
		},
		Items: []domain.Item{
			{
				ChrtId:      1,
				TrackNumber: "MORDOR_TRACK",
				Price:       9999,
				Rid:         "OneRingToRuleThemAll",
				Name:        "The One Ring",
				Sale:        0,
				Size:        "One size fits all Sauron",
				TotalPrice:  9999,
				NmID:        1,
				Brand:       "Sauron Crafts",
				Status:      1,
			},
		}}

	cache.Set(testOrder.OrderUID, &testOrder)
	result, err := cache.GiveFromCache(testOrder.OrderUID)
	require.NoError(t, err, "Посылка не получена")

	assert.Equal(t, "OneRing123", result.OrderUID, "Несоответствие названий")
	assert.Equal(t, "Frodo", result.Delivery.Name, "Неправильное имя получателя")
	assert.Equal(t, "Shir", result.Delivery.City, "Город неправильный")

}
func TestNotfound(t *testing.T) {
	cache := New()

	result, err := cache.GiveFromCache("Не ну такого Uid точно тут не может быть")
	assert.Error(t, err, "Должна быть ошибка")
	assert.Nil(t, result, "Результат должен быть nil")
}

// go test -v ./pkg/cache/ -cover
// для запуска
