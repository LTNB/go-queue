package redis

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"testing"
	"time"
)

var queue *UniversalRedisQueue

func TestSimpleQueue(t *testing.T) {
	fmt.Println("TestSimpleQueue running....")
	dataStr := "data_string"
	universalRedisQueue := UniversalQueueMessage{Data: []byte(dataStr)}
	queue.Queue(universalRedisQueue)
	msg, _ := queue.Take()
	queue.Finish(msg)
}

func TestSimpleRequeue(t *testing.T) {
	fmt.Println("TestSimpleRequeue running....")
	dataStr := "data_string"
	universalRedisQueue := UniversalQueueMessage{Data: []byte(dataStr)}
	queue.Queue(universalRedisQueue)
	message, err := queue.Take()
	assert.Equal(t, message.NumRequeue, 0, "expected requeue is 0")
	if err != nil {
		t.Errorf("expected err is nil...")
	}
	queue.Requeue(message)
	message, err = queue.Take()
	assert.Equal(t, message.NumRequeue, 1, "expected requeue is 0")
	queue.Finish(message)
}

func TestSimpleFinish(t *testing.T) {
	fmt.Println("TestFinish running....")
	dataStr := "data_string"
	universalRedisQueue := UniversalQueueMessage{Data: []byte(dataStr)}
	queue.Queue(universalRedisQueue)
	message, err := queue.Take()
	assert.Equal(t, message.NumRequeue, 0, "expected requeue is 0")
	if err != nil {
		t.Errorf("expected err is nil...")
	}
	queue.Finish(message)
	msg, err := queue.Take()
	assert.NotNil(t, err, nil, "SUCCESS")
	assert.Nil(t, msg.Data, nil, "SUCCESS")
}

func TestPerformanceQueue(t *testing.T) {
	fmt.Println("TestPerformanceQueue running....")
	for k := 0; k < 100; k++ {
		message := UniversalQueueMessage{Data: generateRandomBytes(rand.Intn(100))}
		queue.Queue(message)
	}
	for k := 0; k < 100; k++ {
		msg, _ := queue.Take()
		queue.Finish(msg)
	}
	assert.Equal(t, queue.QueueSize(), 0, "SUCCESS")
}

// queue ==> take ==> getOrphanMessages ==> finish
func TestGetOrphanMessages(t *testing.T) {
	fmt.Println("TestGetOrphanMessages running...")
	for k := 0; k < 10; k++ {
		message := UniversalQueueMessage{Data: generateRandomBytes(rand.Intn(10))}
		queue.Queue(message)
	}
	for k := 0; k < 10; k++ {
		queue.Take()
	}
	time.Sleep(10 * time.Second)
	msg := queue.GetOrphanMessages(10 * 1000)
	assert.Equal(t, len(msg), 10, "Expected receiving 10 message has been queue")
	for k := 0; k < len(msg); k++ {
		queue.Finish(msg[k])
	}
}

func TestQueueSize(t *testing.T) {
	fmt.Println("TestQueueSize running...")
	for k := 0; k < 10; k++ {
		message := UniversalQueueMessage{Data: generateRandomBytes(rand.Intn(10))}
		queue.Queue(message)
	}
	size := queue.QueueSize()
	assert.Equal(t, size, 10, "Expected receiving size = 10 ")
	for k := 0; k < 10; k++ {
		msg, _ := queue.Take()
		queue.Finish(msg)
	}
}

func TestOrphanMessageSize(t *testing.T) {
	fmt.Println("TestOrphanMessageSize running...")
	for k := 0; k < 10; k++ {
		message := UniversalQueueMessage{Data: generateRandomBytes(rand.Intn(10))}
		queue.Queue(message)
	}

	for k := 0; k < 10; k++ {
		queue.Take()
	}
	assert.Equal(t, queue.EphemeralSize(), 10, "Expected receiving size = 10 ")
	msg := queue.GetOrphanMessages(10 * 1000)
	for k := 0; k < len(msg); k++ {
		queue.Finish(msg[k])
	}
}

func generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

func setup() {
	fmt.Println("run before")
	queue = &UniversalRedisQueue{Address: "35.247.157.146:16379", Password: "scte1234", Name: "test_queue"}
	queue.Init()
}

func destroy() {
	queue.Destroy()
}

func TestMain(m *testing.M) {
	setup()
	r := m.Run()
	//destroy()
	os.Exit(r)
}
