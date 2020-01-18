package redis

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

var queueInstance *UniversalRedisQueue

type UniversalRedisQueue struct {
	Address            string
	Password           string
	Name               string
	redisHashName      string
	redisListName      string
	redisSortedSetName string
	PoolSize           int
	client             redis.UniversalClient
}

type UniversalQueueMessage struct {
	queueTimestamp time.Time
	NumRequeue     int
	ID             int64
	Data           []byte
	timestamp      int64
}

func (queue *UniversalRedisQueue) Init() {
	initDefaultConfig(queue)
	queue.client = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{queue.Address},
		Password: queue.Password,
		PoolSize: queue.PoolSize,
	})

	queueInstance = queue
}

func initDefaultConfig(universalRedisQueue *UniversalRedisQueue) {
	if universalRedisQueue.Address == "" {
		universalRedisQueue.Address = "127.0.0.1:6379"
	}
	if universalRedisQueue.Name == "" {
		universalRedisQueue.redisHashName = "default_hash"
		universalRedisQueue.redisSortedSetName = "default_set"
		universalRedisQueue.redisListName = "default_list"
	} else {
		universalRedisQueue.redisHashName = universalRedisQueue.Name + "_hash"
		universalRedisQueue.redisSortedSetName = universalRedisQueue.Name + "_set"
		universalRedisQueue.redisListName = universalRedisQueue.Name + "_list"
	}
	if universalRedisQueue.PoolSize == 0 {
		universalRedisQueue.PoolSize = 10
	}
}

func (queue *UniversalRedisQueue) Destroy() {
	queue.client.Del(queue.redisHashName, queue.redisListName, queue.redisSortedSetName)
}

func (queue *UniversalRedisQueue) GetRedisQueueClient() redis.UniversalClient {
	return queue.client
}

/*
 * put id to messenger_email_l
 * put message to messenger_email_h
 */
func (queue UniversalRedisQueue) Queue(message UniversalQueueMessage) {
	var now = time.Now()
	var nowMillis = now.UnixNano() / int64(time.Millisecond)
	message.ID = nowMillis
	message.NumRequeue = 0
	message.queueTimestamp = now
	message.timestamp = nowMillis
	messageStr, _ := json.Marshal(message)
	field := nowMillis
	queueScript := redis.NewScript(`
		 redis.call("HSET", KEYS[1], ARGV[1], ARGV[2])
		 redis.call("LPUSH", KEYS[2], ARGV[1])
	`)

	queueScript.Run(queueInstance.client, []string{queueInstance.redisHashName, queueInstance.redisListName},
		field, string(messageStr)).Result()
}

/*
 * Pop id from messenger_email_l
 * get message from messenger_email_h
 * push id to messneger_email_s
 */
func (queue UniversalRedisQueue) Take() (UniversalQueueMessage, error) {
	var message UniversalQueueMessage
	queueScript := redis.NewScript(`
			local field = redis.call("RPOP", KEYS[1])
			redis.call("ZADD", KEYS[3], field, field)
			return redis.call("HGET", KEYS[2], field)
			`)
	messageArr, err := queueScript.Run(queueInstance.client, []string{queueInstance.redisListName, queueInstance.redisHashName, queueInstance.redisSortedSetName}).Result()
	if err != nil {
		return message, err
	}
	err = json.Unmarshal([]byte(messageArr.(string)), &message)
	return message, err
}

/*
 * pop id from messeger_email_s
 * push id to messenger_email_l
 * override message in messager_email_h with numRequeue++, reset time in queue
 */
func (queue UniversalRedisQueue) Requeue(message UniversalQueueMessage) {
	var now = time.Now()
	message.NumRequeue = message.NumRequeue + 1
	message.queueTimestamp = now
	messageStr, _ := json.Marshal(message)
	field := strconv.FormatInt(message.ID, 10)
	queueScript := redis.NewScript(`
		redis.call("LPUSH", KEYS[1], ARGV[1])
		redis.call("ZREM", KEYS[2], ARGV[1])
		redis.call("HSET", KEYS[3], ARGV[1], ARGV[2])
		`)
	_, err := queueScript.Run(queueInstance.client, []string{queueInstance.redisListName, queueInstance.redisSortedSetName, queueInstance.redisHashName}, field, messageStr).Result()
	if err != nil {
		fmt.Println("err", err)
	}
}

/*
 * Get id from messeger_email_s not remove
 */
func (queue UniversalRedisQueue) GetOrphanMessages(thresholdTimestampMs int64) []UniversalQueueMessage {
	var result []UniversalQueueMessage
	max := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond)-thresholdTimestampMs, 10)
	dataInQueueStr := queueInstance.client.ZRangeByScore(queueInstance.redisSortedSetName, redis.ZRangeBy{
		Min:    "0",
		Max:    max,
		Offset: 0,
		Count:  100,
	})

	if dataInQueueStr != nil {
		dataInQueue := dataInQueueStr.Val()
		result = make([]UniversalQueueMessage, len(dataInQueue))
		len := len(dataInQueue)
		for i := 0; i < len; i++ {
			var message UniversalQueueMessage
			data, _ := queueInstance.client.HGet(queueInstance.redisHashName, dataInQueue[i]).Bytes()
			json.Unmarshal(data, &message)
			result[i] = message
		}
	}
	return result
}

/*
 * del from hash
 * del from sortedSet
 */
func (queue UniversalRedisQueue) Finish(message UniversalQueueMessage) {
	field := strconv.FormatInt(message.ID, 10)
	queueScript := redis.NewScript(`
		redis.call("HDEL", KEYS[1], ARGV[1])
		redis.call("ZREM", KEYS[2], ARGV[1])
		`)
	queueScript.Run(queueInstance.client, []string{queueInstance.redisHashName, queueInstance.redisSortedSetName}, field).Result()
}

func (queue UniversalRedisQueue) QueueSize() int {
	size := queue.client.LLen(queue.redisListName)
	return int(size.Val())

}

func (queue UniversalRedisQueue) EphemeralSize() int {
	size := queue.client.ZCard(queue.redisSortedSetName)
	return int(size.Val())
}
