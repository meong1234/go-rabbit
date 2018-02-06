package application

import (
	"context"
	"github.com/go-rabbit/amqp"
	"github.com/go-rabbit/messaging"
	"github.com/go-rabbit/util"
)

type (
	subscriber struct {
		queueName string
		broker    messaging.Broker
	}
)

func (app *Application) NewSubscriberDaemon() util.Daemon {
	broker := amqp.NewAmqpBroker(app.rabbit)
	return &subscriber{
		app.queueName,
		broker,
	}
}

func (sub *subscriber) Start() error {
	err := sub.broker.Start()
	if err != nil {
		return err
	}

	_, err = sub.broker.CreateSubscription(sub.queueName, sub.queueName, "", true, 1, sub.handle)
	if err != nil {
		return err
	}

	return nil
}

func (sub *subscriber) Stop() error {
	return sub.broker.Stop()
}

func (sub *subscriber) handle(ctx context.Context, event messaging.Event) error {
	util.SessionLogger(ctx).Debugf("received : %s", string(event.GetBody()))
	return nil
}
