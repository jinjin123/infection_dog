package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"infection/rmq"
	"infection/util/lib"
	"log"
	"os"
	"os/signal"
)

const AmqpURI = "amqp://jin:jinjin123@192.168.50.100:5672"

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	config := rmq.NewIConfigByVHost("infection")
	iConsumer := rmq.NewIConsumerByConfig(AmqpURI, true, false, config)
	//iConsumer.SubscribeToQueue("test", "hostid", false, Handler)
	queuename := lib.HOSTID + "-" + lib.GetRandomString(6)
	cerr := iConsumer.Subscribe("mainupdate", rmq.FanoutExchange, queuename, "hostid", false, Handler)
	if cerr != nil {
		iConsumer.DeleteQueue(queuename)
	}
	s := <-c
	iConsumer.DeleteQueue(queuename)
	fmt.Println("Got signal:", s)

}

type Message struct {
	Hostid string
}

func Handler(d amqp.Delivery) {
	log.Println("接收到消息。。。")
	body := d.Body
	fmt.Println(body)
	consumerTag := d.ConsumerTag
	var msg Message
	json.Unmarshal(body, &msg)
	fmt.Println(msg)
	fmt.Println(consumerTag)
}
