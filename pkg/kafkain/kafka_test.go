package kafkain

import (
	"WB_L0/app/repository/mocks"
	"WB_L0/domain"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

type MockReader struct {
	messages []kafka.Message
	position int
	closed   bool
}

// Моки для тестирования
func NewMockReader(msgs []kafka.Message) *MockReader {
	return &MockReader{
		messages: msgs,
		position: 0,
	}
}

func (m *MockReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	if m.closed {
		return kafka.Message{}, errors.New("Закрытый канал чтения")
	}

	if m.position >= len(m.messages) {
		return kafka.Message{}, errors.New("Сообщения закончились")
	}
	msg := m.messages[m.position]
	m.position++
	return msg, nil
}

func (m *MockReader) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	if m.closed {
		return errors.New("Закрытый канал чтения")
	}
	return nil
}

func (m *MockReader) Close() error {
	m.closed = true
	return nil
}

func TestProcessMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)

	validMsg := kafka.Message{
		Value: []byte(`{"order_uid": "test123", "track_number": "ДолгоПисатьОригинал"}`),
	}
	// Ожидание
	mockRepo.EXPECT().
		SaveOrder(gomock.Any()).
		DoAndReturn(func(order *domain.IncomingData) error {
			assert.Equal(t, "test123", order.OrderUID)
			assert.Equal(t, "ДолгоПисатьОригинал", order.TrackNumber)
			return nil
		})
		// вызов тестируемой функции
	err := ProcessMessage(validMsg)

	assert.NoError(t, err)
}

func TestEmptyJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	empty := kafka.Message{
		Value: []byte(``),
	}

	err := ProcessMessage(empty)

	assert.Error(t, err)
	assert.Equal(t, "Пустой json", err.Error())
}

func TestJsonTroubls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// json  с неправильной структурой
	invalid := kafka.Message{
		Value: []byte(`Тут некорректный`),
	}

	err := ProcessMessage(invalid)

	assert.Error(t, err)
	assert.Equal(t, "Не корректный json", err.Error())
}

func TestStructTroubls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalid := kafka.Message{
		Value: []byte(`"Тут поля":"неправильные"`),
	}

	err := ProcessMessage(invalid)

	assert.Error(t, err)
}

func TestDBErr(t *testing.T) {
	// оригинальная функция
	original := orderRepository
	// откат изменений дефером
	defer func() {
		orderRepository = original
	}()

	orderRepository = func(data *domain.IncomingData) error {
		assert.Equal(t, "test123", data.OrderUID)
		return errors.New("database error")
	}

	msg := kafka.Message{
		Value: []byte(`{"order_uid": "test123", "track_number": "ДолгоПисатьОригинал"}`),
	}

	err := ProcessMessage(msg)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}

func TestDBNoErr(t *testing.T) {
	// оригинальная функция
	original := orderRepository
	// откат изменений дефером
	defer func() {
		orderRepository = original
	}()
	// мок, который успешен
	orderRepository = func(data *domain.IncomingData) error {
		assert.Equal(t, "test123", data.OrderUID)
		assert.Equal(t, "ДолгоПисатьОригинал", data.TrackNumber)
		return nil
	}
	msg := kafka.Message{
		Value: []byte(`{"order_uid": "test123", "track_number": "ДолгоПисатьОригинал"}`),
	}
	err := ProcessMessage(msg)

	assert.NoError(t, err)
}
