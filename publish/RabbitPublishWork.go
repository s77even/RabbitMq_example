package main

import (
	"fmt"
	"imooc_RabbitMQ/RabbitMQ"
	"strconv"
	"time"
)

func main(){
	rabbitmq:= RabbitMQ.NewRabbitMQSimple("imoocSimple")
	rabbitmq.PublishSimple("hello rabbit")
	for i := 0; i <= 100; i++ {
		rabbitmq.PublishSimple("hello"+strconv.Itoa(i))
		fmt.Println(i)
		time.Sleep(time.Second)

	}
}
