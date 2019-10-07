package rmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

//消费者interface
type IConsumer interface {
	Subscribe(exchangeName string, exchangeType string, consumerTag string, handlerFunc func(amqp.Delivery)) error
	SubscribeToQueue(queueName string, consumerTag string, handlerFunc func(amqp.Delivery)) error
	DeleteQueue(queueName string)
}

//消费者interface implement
type Consumer struct {
	conn       *amqp.Connection
	durable    bool
	autoDelete bool
}

//创建一个客户端
//amqpURI rabbitmq url
//默认持久化到硬盘，自动删除无用队列
func NewIConsumer(amqpURI string) *Consumer {
	return NewIConsumerByDurableAndAutoDelete(amqpURI, true, true)
}

//创建一个客户端
//amqpURI rabbitmq url
//durable：是否持久化到硬盘
//autoDelete：当没有连接后是否自动删除
func NewIConsumerByDurableAndAutoDelete(amqpURI string, durable, autoDelete bool) *Consumer {
	if amqpURI == "" {
		panic("Cannot initialize connection to broker, please check your profile. Have you initialized?")
	}

	var err error
	var c Consumer
	c.conn, err = amqp.Dial(fmt.Sprintf("%s/", amqpURI))
	if err != nil {
		panic("Failed to connect to AMQP compatible broker at: " + amqpURI)
	}
	c.durable = durable
	c.autoDelete = autoDelete
	return &c
}

//创建一个客户端
//amqpURI rabbitmq url
//config：带有部分配置参数的消费者
func NewIConsumerByConfig(amqpURI string, durable, autoDelete bool, config IConfig) *Consumer {
	return NewIConsumerByConfigAndDurableAndAutoDelete(amqpURI, true, false, config)
}

//创建一个客户端
//amqpURI rabbitmq url
//durable：是否持久化到硬盘
//autoDelete：当没有连接后是否自动删除
//config：带有部分配置参数的消费者
func NewIConsumerByConfigAndDurableAndAutoDelete(amqpURI string, durable, autoDelete bool,
	config IConfig) *Consumer {
	if amqpURI == "" {
		panic("Cannot initialize connection to broker, please check your profile. Have you initialized?")
	}

	var err error
	var c Consumer
	var conf amqp.Config
	conf.Vhost, conf.ChannelMax, conf.Properties = config.GetConfig()
	c.conn, err = amqp.DialConfig(fmt.Sprintf("%s/", amqpURI), conf)
	if err != nil {
		panic("Failed to connect to AMQP compatible broker at: " + amqpURI)
	}
	c.durable = durable
	c.autoDelete = autoDelete
	return &c
}

//订阅消息
//exchangeName: 交换机名字
//exchangeType：交换机类型
//consumerTag：消费者tag
//autoAck：是否自动ack
//handlerFunc func(delivery amqp.Delivery,autoAck bool)：订阅消息后的业务逻辑处理
func (c *Consumer) Subscribe(exchangeName string, exchangeType string, randomqueue string, consumerTag string, autoAck bool,
	handlerFunc func(delivery amqp.Delivery)) error {
	ch, err := c.conn.Channel()
	// defer ch.Close()
	if err != nil {
		fmt.Printf("%s: %s", "Failed to open a channel", err)
		return err
	}

	err = ch.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		c.durable,    // durable
		c.autoDelete, // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		fmt.Printf("%s: %s", "Failed to register an Exchange", err)
		return err
	}

	log.Printf("declared Exchange, declaring Queue (%s)", "")
	queue, err := ch.QueueDeclare(
		randomqueue,  // name of the queue //random queue bind exchange
		c.durable,    // durable
		c.autoDelete, // delete when usused
		false,        // exclusive
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		fmt.Printf("%s: %s", "Failed to register an Queue", err)
		return err
	}

	log.Printf("declared Queue (%d messages, %d consumers), binding to Exchange (key '%s')",
		queue.Messages, queue.Consumers, exchangeName)

	err = ch.QueueBind(
		queue.Name,   // name of the queue
		consumerTag,  // bindingKey
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	deliveries, err := ch.Consume(
		queue.Name,  // queue
		consumerTag, // consumer
		autoAck,     // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return fmt.Errorf("Failed to register a consumer: %s", err)
	}

	go consumeHandler(deliveries, autoAck, handlerFunc)
	return nil
}

//订阅消息队列中的消息
//queueName：消息队列名字
//consumerTag：消费者tag
//autoAck：是否自动ack
//handlerFunc func(delivery amqp.Delivery,autoAck bool)：订阅消息后的业务逻辑处理
func (c *Consumer) SubscribeToQueue(queueName string, consumerTag string, autoAck bool, handlerFunc func(delivery amqp.Delivery)) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel: %s", err)
	}

	log.Printf("Declaring Queue (%s)", queueName)
	queue, err := ch.QueueDeclare(
		queueName,    // name of the queue
		c.durable,    // durable
		c.autoDelete, // delete when usused
		false,        // exclusive
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("Failed to register an Queue: %s", err)
	}
	deliveries, err := ch.Consume(
		queue.Name,  // queue
		consumerTag, // consumer
		autoAck,     // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return fmt.Errorf("Failed to register a consumer: %s", err)
	}

	go consumeHandler(deliveries, autoAck, handlerFunc)
	return nil
}
func (c *Consumer) DeleteQueue(queueName string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel: %s", err)
	}
	ch.QueueDelete(queueName, false, false, false)
	return nil
}
func consumeHandler(deliveries <-chan amqp.Delivery, autoAck bool, handlerFunc func(d amqp.Delivery)) {
	for d := range deliveries {
		// Invoke the handlerFunc func we passed as parameter.
		handlerFunc(d)
		if !autoAck {
			d.Ack(false)
		}
	}
}
