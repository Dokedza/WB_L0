package kafkain

import (
	"WB_L0/app/repository"
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader ReaderInterfase
	repo   repository.OrderRepository
}
type ReaderInterfase interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type ReaderAdapter struct {
	reader *kafka.Reader
}

func NewReaderAdapter(reader *kafka.Reader) ReaderInterfase {
	return &ReaderAdapter{reader: reader}
}

func (k *ReaderAdapter) Close() error {
	return k.reader.Close()
}
func (k *ReaderAdapter) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	return k.reader.CommitMessages(ctx, msgs...)
}
func (k *ReaderAdapter) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return k.reader.ReadMessage(ctx)
}
