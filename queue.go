package go_queue

import "time"

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

type UniversalMessageQueue struct {
	QueueTimestamp time.Time
	NumRequeue     int
	ID             int64
	Data           []byte
	Timestamp      int64
}

type IQueue interface {
	Queue(message UniversalMessageQueue)
	Take() (UniversalMessageQueue, error)
	Requeue(message UniversalMessageQueue)
	GetOrphanMessages(thresholdTimestampMs int64) []UniversalMessageQueue
	Finish(message UniversalMessageQueue)
	QueueSize() int
	EphemeralSize() int
}
