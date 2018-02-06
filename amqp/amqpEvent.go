package amqp

import "github.com/streadway/amqp"

type amqpEvent struct {
	delivery amqp.Delivery
}

func (e *amqpEvent) GetCorrelationID() string {
	return e.delivery.CorrelationId
}

func (e *amqpEvent) GetBody() []byte {
	return e.delivery.Body
}

func (e *amqpEvent) Ack() error {
	return e.delivery.Ack(false)
}

func (e *amqpEvent) Nack(requeue bool) error {
	return e.delivery.Nack(false, requeue)
}

func (e *amqpEvent) Reject(requeue bool) error {
	return e.delivery.Reject(requeue)
}
