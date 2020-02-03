package pubsub

import (
	"fmt"
	kafka2 "github.com/LTNB/go-queue/kafka"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"testing"
	"time"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

var pubsubInstance *KafkaPubsubService

func setup() {
	fmt.Println("run before")
	configConsumer := &kafka2.UniversalKafkaConsumerConfig{
		BootstrapServers: "localhost:9092",
		GroupId:          "myGroup",
		AutoOffsetReset:  "earliest",
	}

	configConsumer.Init()
	configProducer := &kafka2.UniversalKafkaProducerConfig{
		BootstrapServers: "localhost:9092",
	}
	configProducer.Init()

	c := configConsumer.GetConsumer()
	p := configProducer.GetProducer()
	pubsubInstance = &KafkaPubsubService{
		Producer: p,
		Consumer: c,
		isReady:  false,
	}

	pubsubInstance.Init()
}

func destroy() {
	pubsubInstance.Producer.Close()
	pubsubInstance.Consumer.Close()
}

type SubscriberMock struct {
}

func (mock SubscriberMock) OnMessage(message interface{}) {
	fmt.Println(message)
}

func TestPubsubMessage(t *testing.T) {
	go func() {
		for e := range pubsubInstance.Producer.Events() {
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
		pubsubInstance.Producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}

	// Wait for message deliveries before shutting down
	pubsubInstance.Producer.Flush(15 * 1000)


	pubsubInstance.Subscribe(topic, SubscriberMock{})
	time.Sleep(10 * time.Minute)

}

func TestMain(m *testing.M) {
	setup()
	r := m.Run()
	destroy()
	os.Exit(r)
}
