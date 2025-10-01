package kafkain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	database "WB_L0/app/repository"
	order "WB_L0/domain"

	"github.com/segmentio/kafka-go"
)

// пакетная переменная для зависимостей в тестах
var (
	readerCreate = func() *kafka.Reader {
		return kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{"redpanda:9092"},
			Topic:     "my-topic",
			Partition: 0,
			GroupID:   "my-group",
		})
	}
	massageProcessor = ProcessMessage
	orderRepository  = func(data *order.IncomingData) error {
		return database.DataHasArrivedInOrders(data, nil)
	}
)

// Функция для чтения сообщений из Kafka
func ConsumeMessage() {
	// Создаем новый reader для Kafka
	r := readerCreate()

	defer r.Close()

	for {
		// Читаем сообщение из Kafka
		log.Println("Ожидание сообщений...")
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Не удалось прочесть сообщение:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		err = massageProcessor(msg)
		if err != nil {
			log.Println("Не удалось обработать сообщение:", err)
		}

		// Сообщение успешно обработано, можно проводить коммит
		if err := r.CommitMessages(context.Background(), msg); err != nil {
			log.Println("Не удалось сохранить сообщение:", err)
		}

	}
}

func ProcessMessage(m kafka.Message) error {

	if len(m.Value) == 0 {
		return errors.New("Пустой json")
	}
	if !json.Valid(m.Value) {
		return errors.New("Не корректный json")
	}
	var dataResp order.IncomingData
	err := json.Unmarshal(m.Value, &dataResp)
	if err != nil {
		log.Println("Не удалось декодировать сообщение:", err)
		return err
	}
	fmt.Println("Сообщение успешно обработано")

	return orderRepository(&dataResp)

}
