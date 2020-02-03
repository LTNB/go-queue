package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

type UniversalKafkaConsumerConfig struct {
	BootstrapServers     string //format: host:port,host:port....
	GroupId              string
	SessionTimeOutMs     int             //default 10000
	AutoOffsetReset      string          //earliest, latest, none
	FetchMaxBytes        int             //default 52428800
	FetchMinBytes        int             //default 1
	HeartbeatIntervalMs  int             //must less than timeout default 3000
	ConnectionMaxIdleMs  int             // default 540000
	DefaultApiTimeOutsMs int             //default 60000
	ConfigMap            kafka.ConfigMap //if you want custom default config https://kafka.apache.org/documentation/#consumerconfigs
	consumer             *kafka.Consumer
}

func (consumerConfig *UniversalKafkaConsumerConfig) Init() {
	if consumerConfig.ConfigMap == nil {
		consumerConfig.ConfigMap = kafka.ConfigMap{}
		if consumerConfig.BootstrapServers != "" {
			consumerConfig.ConfigMap["bootstrap.servers"] = consumerConfig.BootstrapServers
		} else {
			consumerConfig.ConfigMap["bootstrap.servers"] = "localhost:9092"
		}
		if consumerConfig.GroupId != "" {
			consumerConfig.ConfigMap["group.id"] = consumerConfig.GroupId
		}
		if consumerConfig.SessionTimeOutMs != 0 {
			consumerConfig.ConfigMap["session.timeout.ms"] = consumerConfig.SessionTimeOutMs
		}
		if consumerConfig.AutoOffsetReset != "" {
			consumerConfig.ConfigMap["auto.offset.reset"] = consumerConfig.AutoOffsetReset
		}
		if consumerConfig.FetchMinBytes != 0 {
			consumerConfig.ConfigMap["fetch.min.bytes"] = consumerConfig.FetchMinBytes
		}
		if consumerConfig.FetchMaxBytes != 0 {
			consumerConfig.ConfigMap["fetch.max.bytes"] = consumerConfig.FetchMaxBytes
		}
		if consumerConfig.HeartbeatIntervalMs != 0 {
			consumerConfig.ConfigMap["heartbeat.interval.ms"] = consumerConfig.HeartbeatIntervalMs
		}
		if consumerConfig.ConnectionMaxIdleMs != 0 {
			consumerConfig.ConfigMap["connections.max.idle.ms"] = consumerConfig.ConnectionMaxIdleMs
		}
		if consumerConfig.DefaultApiTimeOutsMs != 0 {
			consumerConfig.ConfigMap["default.api.timeout.ms"] = consumerConfig.DefaultApiTimeOutsMs
		}
	}
	c, err := kafka.NewConsumer(&consumerConfig.ConfigMap)

	if err != nil {
		log.Panic(err)
	}

	consumerConfig.consumer = c
}

func (consumerConfig *UniversalKafkaConsumerConfig) GetConsumer() *kafka.Consumer {
	return consumerConfig.consumer
}

type UniversalKafkaProducerConfig struct {
	BootstrapServers    string          //format: host:port,host:port....
	Acks                string          //default 1
	BufferMemory        int             //default 33554432
	CompressionType     string          // default none
	Retries             int             //default 2147483647
	ConnectionMaxIdleMs int             // default 540000
	ConfigMap           kafka.ConfigMap //if you want custom default config https://kafka.apache.org/documentation/#producerconfigs
	producer            *kafka.Producer
}

func (producerConfig *UniversalKafkaProducerConfig) Init() {
	if producerConfig.ConfigMap == nil {
		producerConfig.ConfigMap = kafka.ConfigMap{}
		if producerConfig.BootstrapServers != "" {
			producerConfig.ConfigMap["bootstrap.servers"] = producerConfig.BootstrapServers
		}
		if producerConfig.Acks != "" {
			producerConfig.ConfigMap["acks"] = producerConfig.Acks
		}
		if producerConfig.BufferMemory != 0 {
			producerConfig.ConfigMap["buffer.memory"] = producerConfig.BufferMemory
		}
		if producerConfig.CompressionType != "" {
			producerConfig.ConfigMap["compression.type"] = producerConfig.CompressionType
		}
		if producerConfig.Retries != 0 {
			producerConfig.ConfigMap["retries"] = producerConfig.Retries
		}
		if producerConfig.ConnectionMaxIdleMs != 0 {
			producerConfig.ConfigMap["connection.max.idle.ms"] = producerConfig.ConnectionMaxIdleMs
		}
	}
	p, err := kafka.NewProducer(&producerConfig.ConfigMap)
	if err != nil {
		log.Panic(err)
	}
	producerConfig.producer = p
}

func (producerConfig *UniversalKafkaProducerConfig) GetProducer() *kafka.Producer {
	return producerConfig.producer
}
