package pubsub

import (
	"encoding/json"
	"fmt"
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
			Address: "35.247.157.146:16379", Password: "scte1234",
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

func (mock SubscriberPubSubMessageMock) OnMessage(message interface{}) {
	fmt.Println(message)
	//result := go_queue.UniversalPubSubMessage{}
	//countUniversalMessage++
	//if err := json.Unmarshal([]byte(message), &result); err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(result)
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
	subscriber1 := SubscriberPubSubMessageMock{}
	pubsubInstace.Subscribe("localhost", subscriber)
	pubsubInstace.Subscribe("localhost", subscriber1)
	pubsubInstace.Publish("localhost", msg)
	time.Sleep(15 * time.Second)
	assert.NotNil(t, countUniversalMessage)
}

type SubscriberCommonsMessageMock struct {
}

var countCommonsMessage = 0

func (subscriber SubscriberCommonsMessageMock) OnMessage(message interface{}) {
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
