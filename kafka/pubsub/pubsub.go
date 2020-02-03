package pubsub

import (
	"encoding/json"
	"fmt"
	go_queue "github.com/LTNB/go-queue"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"sync"
	"time"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

type KafkaPubsubService struct {
	Producer   *kafka.Producer
	Consumer   *kafka.Consumer
	subscriber map[string]map[go_queue.ISubscriber]bool
	isReady    bool
	l sync.Mutex
}

func (psService *KafkaPubsubService) Init() {
	initDefaultConfig(psService)
	if psService.Consumer == nil {
		return
	}
	go func() {
		psService.Consumer.SubscribeTopics(psService.getTopic(), nil) // map from subscriber to string topics
		for {
			if psService.isReady {
				msg, err := psService.Consumer.ReadMessage(-1)
				if err != nil {
					fmt.Println(err)
					continue
				}
				for sub, _ := range psService.subscriber[*msg.TopicPartition.Topic] {
					sub.OnMessage(string(msg.Value))
				}
			} else {
				psService.Consumer.SubscribeTopics(psService.getTopic(), nil)
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

func (psService *KafkaPubsubService) getTopic() []string {
	topic := make([]string, 0)
	for k, _ := range psService.subscriber {
		topic = append(topic, k)
	}
	return topic
}

func initDefaultConfig(redisPubSubService *KafkaPubsubService) {
	if redisPubSubService.subscriber == nil {
		redisPubSubService.subscriber = make(map[string]map[go_queue.ISubscriber]bool)
	}
}

func (psService KafkaPubsubService) CreateMessage() go_queue.UniversalPubSubMessage {
	return go_queue.UniversalPubSubMessage{
		ID:   time.Now().UnixNano() / int64(time.Millisecond),
		Data: nil,
	}
}

func (psService KafkaPubsubService) CreateMessageWithData(data []byte) go_queue.UniversalPubSubMessage {
	return go_queue.UniversalPubSubMessage{
		ID:   time.Now().UnixNano() / int64(time.Millisecond),
		Data: data,
	}
}

func (psService KafkaPubsubService) PublishUniversal(topic string, message go_queue.UniversalPubSubMessage) (bool, error) {
	msgByte, err := json.Marshal(message)
	if err != nil {
		return false, err
	}
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msgByte,
	}
	if err = psService.Producer.Produce(kafkaMsg, nil); err != nil {
		return false, err
	}
	return true, nil
}

func (psService KafkaPubsubService) Publish(topic string, message interface{}) (bool, error) {
	msgByte, err := json.Marshal(message)
	if err != nil {
		return false, err
	}
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msgByte,
	}
	if err = psService.Producer.Produce(kafkaMsg, nil); err != nil {
		return false, err
	}
	return true, nil
}

func (psService *KafkaPubsubService) Subscribe(channel string, subscriber go_queue.ISubscriber) error {
	psService.isReady = false
	psService.l.Lock()
	defer psService.l.Unlock()
	if psService.subscriber[channel] == nil {
		psService.subscriber[channel] = make(map[go_queue.ISubscriber]bool)
	}
	psService.subscriber[channel][subscriber] = true
	psService.isReady = true
	return nil
}

func (psService *KafkaPubsubService) Unsubscribe(channel string, subscriber go_queue.ISubscriber) error {
	psService.isReady = false
	psService.l.Lock()
	defer psService.l.Unlock()
	delete(psService.subscriber[channel], subscriber)
	psService.isReady = true
	return nil
}
