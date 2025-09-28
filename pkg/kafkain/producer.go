package kafkain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	database "WB_L0/app/repository"
	order "WB_L0/domain"

	"github.com/segmentio/kafka-go"
)

// Функция для чтения сообщений из Kafka
func ConsumeMessage() {
	topic := "my-topic"
	partition := 0
	groupID := "my-group"
	// Создаем новый reader для Kafka
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"redpanda:9092"},
		Topic:     topic,
		Partition: partition,
		GroupID:   groupID,
	})

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

		var dataResp order.IncomingData
		err = json.Unmarshal(msg.Value, &dataResp)
		if err != nil {
			log.Println("Не удалось декодировать сообщение:", err)
			continue
		}

		// Обрабатываем сообщение
		err = database.DataHasArrivedInOrders(&dataResp, nil)
		if err != nil {
			log.Println("Не удалось обработать сообщение:", err)
			continue
		}

		// Сообщение успешно обработано, можно проводить коммит
		if err := r.CommitMessages(context.Background(), msg); err != nil {
			log.Println("Не удалось сохранить сообщение:", err)
		}

		fmt.Println("Сообщение успешно обработано")
	}
}
