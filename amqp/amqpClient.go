package amqp

import (
	"github.com/go-rabbit/util"
	"github.com/streadway/amqp"
)

type amqpClient struct {
	amqpPublisherManager
	amqpSubscriptionManager
	connection *amqp.Connection
}

func (cli *amqpClient) Init(connection *amqp.Connection) error {
	cli.connection = connection

	util.Log.Infof("initiate publisher manager")
	if err := cli.amqpPublisherManager.Init(connection); err != nil {
		util.Log.Warnf("Fail initiate publisher manager %v", err)
		return err
	}

	util.Log.Infof("initiate subscription manager")
	if err := cli.amqpSubscriptionManager.Init(connection); err != nil {
		util.Log.Warnf("Fail initiate subscription manager %v", err)
		return err
	}
	return nil
}

func (cli *amqpClient) Close() error {
	util.Log.Infof("try close subscription manager")
	if err := cli.amqpSubscriptionManager.Close(); err != nil {
		util.Log.Errorf("Failed to close subscription manager: %v\n", err)
	}

	util.Log.Infof("try close publisher manager")
	if err := cli.amqpPublisherManager.Close(); err != nil {
		util.Log.Errorf("Failed to close publisher manager: %v\n", err)
	}

	util.Log.Infof("try close connection")
	return cli.connection.Close()
}
