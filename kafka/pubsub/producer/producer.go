package main

import (
	"fmt"
	kafka2 "github.com/LTNB/go-queue/kafka"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

func main() {
	config := &kafka2.UniversalKafkaProducerConfig{
		BootstrapServers: "localhost:9092",
	}
	config.Init()
	p := config.GetProducer()
	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "test"
	for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "client"} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}
