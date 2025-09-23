package kafkain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	database "github.com/dokedza/WB_L0/app/repository"
	order "github.com/dokedza/WB_L0/domain"
	"github.com/segmentio/kafka-go"
)

// Функция для чтения сообщений из Kafka
func ConsumeMessage() {
	topic := "my-topic"
	partition := 0
	groupID := "my-group"
	// Создаем новый reader для Kafka
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:29092"},
		Topic:     topic,
		Partition: partition,
		GroupID:   groupID,
	})

	defer r.Close()

	for {
		// Читаем сообщение из Kafka
		log.Println("Waiting for new messages...")
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("failed to read message:", err)
		}

		var dataResp order.IncomingData
		err = json.Unmarshal(msg.Value, &dataResp)
		if err != nil {
			log.Println("failed to unmarshal message:", err)
			continue
		}

		// Обрабатываем сообщение
		err = database.DataHasArrivedInOrders(&dataResp, nil)
		if err != nil {
			log.Println("failed to process message:", err)
			continue
		}

		// Сообщение успешно обработано, можно проводить коммит
		if err := r.CommitMessages(context.Background(), msg); err != nil {
			log.Println("failed to commit message:", err)
		}

		fmt.Println("Message processed successfully")
	}
}
