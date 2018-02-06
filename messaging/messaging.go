package messaging

import "context"

type (
	Publisher interface {
		Init(option map[string]interface{}) error
		Publish(payload interface{}, cid string) error
	}

	Event interface {
		GetCorrelationID() string
		GetBody() []byte
		Ack() error
		Nack(requeue bool) error
		Reject(requeue bool) error
	}

	Subscription interface {
		Close() error
	}

	SubscribeProcessor = func(ctx context.Context, event Event) error

	Broker interface {
		Start() error
		Stop() error
		CreateSubscription(queueName, key, exchange string, autoAck bool, workerCount int, processor SubscribeProcessor) (Subscription, error)
		CreatePublisher(publishKey string) (Publisher, error)
	}
)
