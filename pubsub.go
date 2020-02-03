package go_queue

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */
// create standard message
// support some method for business
type UniversalPubSubMessage struct {
	ID   int64
	Data []byte
}

type PubSub interface {
	CreateMessage() UniversalPubSubMessage

	CreateMessageWithData(data []byte) UniversalPubSubMessage

	PublishUniversal(channel string, message UniversalPubSubMessage) (bool, error)

	Publish(channel string, message interface{}) (bool, error)

	Subscribe(channel string, subscriber ISubscriber) error

	Unsubscribe(channel string, subscriber ISubscriber) error
}

//implement for receiving message
type ISubscriber interface {
	OnMessage(message interface{})
}