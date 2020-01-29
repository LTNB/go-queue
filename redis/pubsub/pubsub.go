package pubsub

import (
	"encoding/json"
	go_queue "github.com/LTNB/go-queue"
	conf "github.com/LTNB/go-queue/redis"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */
var l sync.Mutex
var pubsubInstace *RedisPubSubService

type RedisPubSubService struct {
	conf.UniversalRedisConfig
	client     redis.UniversalClient
	Pattern    string
	subscriber map[string]map[go_queue.ISubscriber]bool
	isReady    bool
}

func (psService *RedisPubSubService) Init() {
	initDefaultConfig(psService)
	psService.client = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{psService.Address},
		Password: psService.Password,
		PoolSize: psService.PoolSize,
	})

	pubsubInstace = psService
	go func() {
		pubSub := psService.client.PSubscribe(psService.Pattern)
		for {
			if psService.isReady {
				msg, _ := pubSub.ReceiveMessage()
				go func() {
					for sub, _ := range psService.subscriber[msg.Channel] {
						sub.OnMessage(msg.Payload)
					}
				}()
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

func initDefaultConfig(redisPubSubService *RedisPubSubService) {
	if redisPubSubService.subscriber == nil {
		redisPubSubService.subscriber = make(map[string]map[go_queue.ISubscriber]bool)
	}
	if redisPubSubService.Address == "" {
		redisPubSubService.Address = "127.0.0.1:6379"
	}
	if redisPubSubService.PoolSize == 0 {
		redisPubSubService.PoolSize = 10
	}
}

func (psService RedisPubSubService) CreateMessage() go_queue.UniversalPubSubMessage {
	return go_queue.UniversalPubSubMessage{
		ID:   time.Now().UnixNano() / int64(time.Millisecond),
		Data: nil,
	}
}

func (psService RedisPubSubService) CreateMessageWithData(data []byte) go_queue.UniversalPubSubMessage {
	return go_queue.UniversalPubSubMessage{
		ID:   time.Now().UnixNano() / int64(time.Millisecond),
		Data: data,
	}
}

func (psService RedisPubSubService) Publish(channel string, message go_queue.UniversalPubSubMessage) (bool, error) {
	msgJson, err := json.Marshal(message)
	if err != nil {
		return false, err
	}
	result := psService.client.Publish(channel, msgJson)
	if result.Val() > 0 {
		return true, nil
	}
	return false, result.Err()
}

func (psService *RedisPubSubService) Subscribe(channel string, subscriber go_queue.ISubscriber) error {
	psService.isReady = false
	l.Lock()
	defer l.Unlock()
	if psService.subscriber[channel] == nil {
		psService.subscriber[channel] = make(map[go_queue.ISubscriber]bool)
	}
	psService.subscriber[channel][subscriber] = true
	psService.isReady = true
	return nil
}

func (psService *RedisPubSubService) Unsubscribe(channel string, subscriber go_queue.ISubscriber) error {
	psService.isReady = false
	l.Lock()
	defer l.Unlock()
	delete(psService.subscriber[channel], subscriber)
	psService.isReady = true
	return nil
}

func (psService RedisPubSubService) CountSubscriber() int {
	count := 0
	for _, subMap := range psService.subscriber {
		count += len(subMap)
	}
	return count
}

func (psService RedisPubSubService) GetInstance() RedisPubSubService {
	return *pubsubInstace
}
