package pubsub

import (
	"fmt"
	kafka2 "github.com/LTNB/go-queue/kafka"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
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
		GroupId:          "testGroup",
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

	// Produce messages to topic (asynchronously)
	topic := "test"
	for _, word := range []string{"message", "from", "test", "topic"} {
		pubsubInstance.Producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}

	// Wait for message deliveries before shutting down
	pubsubInstance.Producer.Flush(15 * 1000)


	err := pubsubInstance.Subscribe(topic, SubscriberMock{})
	time.Sleep(10 * time.Second)
	assert.Nil(t, err)

}

func TestMain(m *testing.M) {
	setup()
	r := m.Run()
	destroy()
	os.Exit(r)
}
