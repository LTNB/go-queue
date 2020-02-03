# go-queue
#### Utility and implement base Queue/Pubsub.
#### Init and maintain by LamTNB (baolam0307@gmail.com)
## Installation
Latest release version: `0.1.0`. See [RELEASE-NOTES.md](RELEASE-NOTES.md).
## Usage
`import github.com/go-redis/redis`
`import github.com/confluentinc/confluent-kafka-go/kafka`

## kafka test
1. `cd kafka/kafka_2.12-2.4.0`
2. `bin/zookeeper-server-start.sh config/zookeeper.properties`
3. `bin/kafka-server-start.sh config/server.properties`
4. Run all test case in `kafka/pubsub/pubsub_test.go`
## redis test
1. Setup redis on docker, refer: `https://hub.docker.com/_/redis`
2. Run all test case in `redis/pubsub/pubsub_test.go` for pubsub testing
3. Run all test case in `redis/queue/queue_test.go` for queue testing
