package main

import (
	"fmt"
	kafka2 "github.com/LTNB/go-queue/kafka"
	"github.com/LTNB/go-queue/kafka/pubsub"
	"time"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

type KafkaSubscriberMock struct {

}

func (subscriber KafkaSubscriberMock) OnMessage(message string){
	fmt.Println(message)
}

func main() {
	config := &kafka2.UniversalKafkaConsumerConfig{
		BootstrapServers: "localhost:9092",
		GroupId:          "myGroup",
		AutoOffsetReset:  "earliest",
	}

	config.Init()
	configProducer := &kafka2.UniversalKafkaProducerConfig{
		BootstrapServers: "localhost:9092",
	}
	config.Init()
	p := configProducer.GetProducer()
	psService := pubsub.KafkaPubsubService{
		Producer: p,
		Consumer: config.GetConsumer(),
	}
	psService.Init()

	psService.Subscribe("test", KafkaSubscriberMock{})

	time.Sleep(10*time.Minute)
	//c.SubscribeTopics([]string{"test"}, nil)
	//
	//for {
	//	msg, err := c.ReadMessage(-1)
	//	if err == nil {
	//		fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
	//	} else {
	//		// The client will automatically try to recover from all errors.
	//		fmt.Printf("Consumer error: %v (%v)\n", err, msg)
	//	}
	//}
	//
	//c.Close()
}