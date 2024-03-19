package utils

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var KafkaProducer *kafka.Producer

func SendMetadata(Topic string, metadata RecordPayload, producer *kafka.Producer) (string, error) {

	// Convert struct to bytes
	postBytes, err := json.Marshal(metadata)
	if err != nil {
		log.Fatal(err)
	}

	deliveryChan := make(chan kafka.Event)
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &Topic, Partition: kafka.PartitionAny},
		Value:          postBytes,
	}, deliveryChan)

	if err != nil {
		msg := fmt.Sprintf("Failed to produce message 1: %v\n", err)

		return msg, err

	} else {

		// Wait for delivery report
		e := <-deliveryChan
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			msg := fmt.Sprintf("Delivery failed: %v\n", m.TopicPartition.Error)

			return msg, m.TopicPartition.Error

		} else {
			fmt.Printf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		}
	}

	return "Message Delivered Successfully", nil
}

func InitializeProducer() (*kafka.Producer, error) {

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "test",
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return nil, err
	}

	fmt.Println("Connected to Apache Kafka!")

	return producer, nil
}
