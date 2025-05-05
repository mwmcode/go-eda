package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var kafkaTopic = "my-topic"

func main() {
	publisher, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})

	if err != nil {
		log.Fatalf("Failed to create producer %s", err)
	}
	defer publisher.Close()

	go func() {
		for e := range publisher.Events() {
			switch evt := e.(type) {
			case *kafka.Message:
				if evt.TopicPartition.Error != nil {
					fmt.Printf("failed to deliver message %v\n", evt.TopicPartition)
				} else {
					fmt.Printf("success! produced record to topic %s partition [%d] offset %v\n",
						*evt.TopicPartition.Topic,
						evt.TopicPartition.Partition,
						evt.TopicPartition.Offset,
					)
				}
			}
		}
	}()

	messageTexts := []string{
		`{"orderId": 11, "status": "success"}`, `{"orderId": 45, "status": "failed"}`,
		`{"orderId": 67, "status": "pending"}`, `{"orderId": 90, "status": "initial"}`,
	}

	for i := 0; i < len(messageTexts); i++ {
		err := publisher.Produce(
			&kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &kafkaTopic,
					Partition: kafka.PartitionAny,
				},
				Value: []byte(messageTexts[i]),
			},
			nil,
		)

		if err != nil {
			log.Printf("failed to publish message %v", err)
		}
	}

	publisher.Flush(15 * 1000)
}
