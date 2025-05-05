package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var kafkaTopic = "my-topic"

func main() {
	subscriber, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("failed to create subscriber %s", err)
	}
	defer subscriber.Close()

	err = subscriber.SubscribeTopics([]string{kafkaTopic}, nil)
	if err != nil {
		log.Fatalf("failed to subscribe to topic %s", err)
	}

	fmt.Println("subscriber started, waiting on messages...")

	// poll
	for {
		msg, err := subscriber.ReadMessage(-1)
		if err != nil {
			fmt.Printf("subscriber error: %v (%v)\n", err, msg)
			continue
		}
		fmt.Printf("consumed message: %s: %s\n", msg.TopicPartition, string(msg.Value))
	}
}
