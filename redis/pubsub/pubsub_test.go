package pubsub

import (
	"encoding/json"
	"fmt"
	go_queue "github.com/LTNB/go-queue"
	"github.com/LTNB/go-queue/redis"
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
var pubsubInstance *RedisPubSubService

func setup() {
	fmt.Println("run before")
	pubsubInstace = &RedisPubSubService{
		UniversalRedisConfig: redis.UniversalRedisConfig{
			Address: "localost:6379", Password: "12345",
		},
		Pattern: "*",
	}
	pubsubInstace.Init()
	pubsubInstance = pubsubInstace
}

func TestGetInstance(t *testing.T) {
	redisPubsubInstance := pubsubInstace.GetInstance()
	assert.NotNil(t, redisPubsubInstance, "Create success")
}

func TestCreateMessage(t *testing.T) {
	msg := pubsubInstace.CreateMessage()
	assert.NotNil(t, msg, "Create success")
}

func TestCreateMessageWithData(t *testing.T) {
	msgData := "data mock"
	msgDataByte, _ := json.Marshal(msgData)
	msg := pubsubInstace.CreateMessageWithData(msgDataByte)
	assert.NotNil(t, msg.ID, "Create success")
}

func TestPublish(t *testing.T) {
	msgData := "data mock"
	msgDataByte, _ := json.Marshal(msgData)
	msg := pubsubInstace.CreateMessageWithData(msgDataByte)
	result, err := pubsubInstace.Publish("localhost", msg)
	assert.Nil(t, err, "Expected error is nil")
	assert.True(t, true, result, "Expected publish message success")
}

type SubscriberPubSubMessageMock struct {
}

var countUniversalMessage = 0

func (mock SubscriberPubSubMessageMock) OnMessage(message string) {
	result := go_queue.UniversalPubSubMessage{}
	countUniversalMessage++
	if err := json.Unmarshal([]byte(message), &result); err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}


func TestUnsubscribe(t *testing.T) {
	subscriber := SubscriberPubSubMessageMock{}
	pubsubInstace.Subscribe("new_channel", subscriber)
	assert.Equal(t, 1, pubsubInstace.CountSubscriber(), "expected registry subscriber success")
	pubsubInstace.Unsubscribe("new_channel", subscriber)
	assert.Equal(t, 0, pubsubInstace.CountSubscriber(), "expected remove subscriber success")
}

func TestGetUniversalPubsubMessage(t *testing.T) {
	msgData := "data mock"
	msgDataByte, _ := json.Marshal(msgData)
	msg := pubsubInstace.CreateMessageWithData(msgDataByte)
	subscriber := SubscriberPubSubMessageMock{}
	pubsubInstace.Subscribe("localhost", subscriber)
	pubsubInstace.Publish("localhost", msg)
	time.Sleep(2 * time.Second)
	assert.NotNil(t, countUniversalMessage)
}

type SubscriberCommonsMessageMock struct {
}

var countCommonsMessage = 0

func (subscriber SubscriberCommonsMessageMock) OnMessage(message string) {
	countCommonsMessage++
	fmt.Println(message)

}

func TestGetCommonsMessage(t *testing.T) {
	pubsubInstace = &RedisPubSubService{
		UniversalRedisConfig: redis.UniversalRedisConfig{
			Address: "localhost:6378",
		},
		Pattern: "gps*",
	}
	pubsubInstace.Init()
	pubsubInstance = pubsubInstace

	subscriber := SubscriberCommonsMessageMock{}
	pubsubInstace.Subscribe("tracking", subscriber)
	time.Sleep(3 * time.Second)
	assert.NotNil(t, countCommonsMessage)
}

func TestMain(m *testing.M) {
	setup()
	r := m.Run()
	//destroy()
	os.Exit(r)
}
