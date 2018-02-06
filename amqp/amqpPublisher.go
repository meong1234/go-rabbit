package amqp

import (
	"encoding/json"
	"fmt"
	"github.com/go-rabbit/messaging"
	"github.com/go-rabbit/util"
	"github.com/streadway/amqp"
)

type (
	amqpPublisherManager struct {
		connection *amqp.Connection
		Publishers map[string]*amqpPublisher
	}

	amqpPublisher struct {
		publishKey     string
		publishChannel *amqp.Channel
	}
)

const (
	messagingContentType = "application/json"
)

func (pub *amqpPublisherManager) Init(connection *amqp.Connection) error {
	pub.connection = connection
	var err error
	if pub.Publishers != nil && len(pub.Publishers) > 0 {
		for _, publisher := range pub.Publishers {
			err = publisher.init(connection)
		}
	}
	return err
}

func (pub *amqpPublisherManager) CreatePublisher(publishKey string) (messaging.Publisher, error) {
	util.Log.WithField("publishKey", publishKey).Info("create publisher")
	publisher := amqpPublisher{
		publishKey: publishKey,
	}

	if err := publisher.init(pub.connection); err != nil {
		return nil, err
	}

	if len(pub.Publishers) == 0 {
		pub.Publishers = make(map[string]*amqpPublisher)
	}

	pub.Publishers[publishKey] = &publisher
	return &publisher, nil
}

func (pub *amqpPublisherManager) Close() error {
	var err error
	for _, publisher := range pub.Publishers {
		err = publisher.Close()
	}
	return err
}

func (client *amqpPublisher) init(connection *amqp.Connection) error {
	util.Log.Infof("Amqp publisher initiated %s", client.publishKey)
	channel, err := connection.Channel()
	if err != nil {
		return err
	}
	client.publishChannel = channel

	return nil
}

func (client *amqpPublisher) Publish(payload interface{}, cid string) error {
	if client.publishChannel == nil {
		return fmt.Errorf("Connection currently not available.")
	}

	exchangeName := ""
	routingKey := client.publishKey

	data, err := json.Marshal(payload)
	if err != nil {
		// Failed to encode payload
		return err
	}

	publishing := amqp.Publishing{
		ContentType:     messagingContentType,
		ContentEncoding: "UTF-8",
		CorrelationId:   cid,
		DeliveryMode:    amqp.Persistent,
		Body:            data,
	}

	return client.publishChannel.Publish(exchangeName, routingKey, false, false, publishing)
}

func (client *amqpPublisher) Init(option map[string]interface{}) error {
	if _, err := client.publishChannel.QueueDeclare(client.publishKey, true, false, false, false, option); err != nil {
		return err
	}

	return nil
}

func (client *amqpPublisher) Close() error {
	return client.publishChannel.Close()
}
