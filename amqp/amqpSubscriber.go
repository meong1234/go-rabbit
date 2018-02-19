package amqp

import (
	"fmt"
	"github.com/go-rabbit/messaging"
	"github.com/go-rabbit/util"
	"github.com/streadway/amqp"
	"math/rand"
	"os"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func getConsumerID(queue string) string {
	var hostname string

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	pid := os.Getpid()

	return fmt.Sprintf("%s#%d@%s/%d", queue, pid, hostname, rand.Uint32())
}

func newAmqpSubscription(queueName, key, exchange string, autoAck bool, workerCount int, processor messaging.SubscribeProcessor) *amqpSubscriber {
	return &amqpSubscriber{
		QueueName:   queueName,
		Key:         key,
		Exchange:    exchange,
		AutoAck:     autoAck,
		workerCount: workerCount,
		processor:   processor,
		Logger:      util.Log.WithFields(util.BuildContext("amqpSubscription")).WithField("queueName", queueName),
	}
}

type amqpSubscriptionManager struct {
	connection  *amqp.Connection
	Subscribers []*amqpSubscriber
}

func (sub *amqpSubscriptionManager) Init(connection *amqp.Connection) error {
	sub.connection = connection

	if sub.Subscribers != nil && len(sub.Subscribers) > 0 {
		for _, subscriber := range sub.Subscribers {
			subscriber.init(connection)
		}
	}
	return nil
}

func (sub *amqpSubscriptionManager) CreateSubscription(queueName, key, exchange string, autoAck bool, workerCount int, processor messaging.SubscribeProcessor) (messaging.Subscription, error) {
	subscription := newAmqpSubscription(queueName, key, exchange, autoAck, workerCount, processor)
	if err := subscription.init(sub.connection); err != nil {
		return nil, err
	}
	sub.Subscribers = append(sub.Subscribers, subscription)
	return subscription, nil
}

func (client *amqpSubscriptionManager) Close() error {
	var err error
	for _, sub := range client.Subscribers {
		err = sub.Close()
	}
	return err
}

type amqpSubscriber struct {
	AutoAck      bool
	QueueName    string
	Key          string
	Exchange     string
	consumerId   string
	workerCount  int
	processingWG sync.WaitGroup // use wait group to make sure task processing completes on interrupt signal
	processor    messaging.SubscribeProcessor
	stopChan     chan int
	channel      *amqp.Channel
	Logger       util.Logger
}

func (c *amqpSubscriber) init(connection *amqp.Connection) error {
	c.consumerId = getConsumerID(c.QueueName)
	c.Logger.Infof("initiated")

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	c.stopChan = make(chan int)
	errorChan := channel.NotifyClose(make(chan *amqp.Error))

	go func() {
		// NOTE: We need this in case the channel gets broken (but not the connection), e.g.
		// when an invalid message is acknowledged.
		for err := range errorChan {
			c.Logger.Errorf("queue-channel error %v\n", err)
			c.channel.Cancel(c.consumerId, false)
			c.stopChan <- 1
			//if cerr := c.init(connection); cerr != nil {
			//	util.Log.Errorf("Failed to reestablish channel for queue %s: %v\n", c.QueueName, cerr)
			//}
		}
	}()

	if err := channel.Qos(c.workerCount*2, 0, false); err != nil {
		return err
	}

	deadLetterQueue := c.QueueName + "_DEAD_LETTER"
	if _, err := channel.QueueDeclare(deadLetterQueue, true, false, false, false, nil); err != nil {
		return err
	}

	table := make(amqp.Table)
	table["x-dead-letter-exchange"] = ""
	table["x-dead-letter-routing-key"] = deadLetterQueue
	if _, err := channel.QueueDeclare(c.QueueName, true, false, false, false, table); err != nil {
		return err
	}

	if (c.Key != "") && (c.Exchange != "") {
		if err := channel.QueueBind(c.QueueName, c.Key, c.Exchange, false, nil); err != nil {
			return err
		}
	}

	deliveryChan, err := channel.Consume(c.QueueName, c.consumerId, c.AutoAck, false, false, false, nil)
	if err != nil {
		return err
	}

	c.channel = channel

	go c.runEventProcessor(deliveryChan)

	c.Logger.Infof("initiated done")
	return nil
}

func (c *amqpSubscriber) runEventProcessor(deliveryChan <-chan amqp.Delivery) {
	pool := util.NewIntPool(c.workerCount)

	for {
		select {
		case d := <-deliveryChan:
			if len(d.Body) == 0 {
				d.Nack(false, false)
				return
			}

			c.Logger.Debugf("recieved message %s\n", d.CorrelationId)

			workerNumber := pool.Get()
			event := amqpEvent{d}
			c.processingWG.Add(1)

			go func(worker int, event *amqpEvent, log util.Logger) {
				logger := log.WithField("workerNumber", worker).
					WithField("correlationId", event.GetCorrelationID())

				ctx := util.NewSessionCtx(event.GetCorrelationID(), logger)

				logger.WithField("message", string(event.GetBody())).Debug("process event started")
				c.processor(ctx, event)
				logger.Debug("process event finish")

				defer func() {
					c.processingWG.Done()
					pool.Put(worker)
				}()

			}(workerNumber, &event, c.Logger)

		case <-c.stopChan:
			c.Logger.Debugf("stop triggered")
			return
		}
	}
}

func (c *amqpSubscriber) Close() error {
	c.Logger.Debugf("close triggered")
	close(c.stopChan)

	// Waiting for any tasks being processed to finish
	c.processingWG.Wait()

	// NOTE: This should close the deliveryChannel, which quits the loop in Run(), which stops this subscriber
	if err := c.channel.Cancel(c.consumerId, false); err != nil {
		return err
	}
	return c.channel.Close()
}
