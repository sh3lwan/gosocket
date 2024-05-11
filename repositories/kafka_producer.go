package repositories

import (
	"encoding/json"
	"fmt"
	"log"

	. "github.com/sh3lwan/gosocket/types"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	TOPIC = "notifications"
	HOST  = "kafka"
	PORT  = "9092"
)

type Producer struct {
	*kafka.Producer
}

func NewProducer() *Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", HOST, PORT),
	})

	if err != nil {
		log.Fatal(err.Error())

	}
	return &Producer{p}
}

func (p *Producer) WriteNotification(message ReceivedMessage) error {
	topic := TOPIC

	msg, err := json.Marshal(message)

	if err != nil {
		return err
	}

	return p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil)
}
