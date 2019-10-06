package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"infection/rmq"
	"log"
	"net/http"
)

const AmqpURI = "amqp://jin:jinjin123@192.168.50.100:5672"

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

	})
	config := rmq.NewIConfigByVHost("infection")
	iConsumer := rmq.NewIConsumerByConfig(AmqpURI, true, false, config)
	iConsumer.SubscribeToQueue("clearpic", "hostid", false, Handler)
	//iConsumer.Subscribe("killip",rmq.DirectExchange,"hostid",false,Handler)
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		panic(err)
	}
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
