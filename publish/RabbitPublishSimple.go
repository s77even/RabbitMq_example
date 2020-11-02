package main

import (
	"fmt"
	"imooc_RabbitMQ/RabbitMQ"
)

func main(){
	rabbitmq:= RabbitMQ.NewRabbitMQSimple("imoocSimple")
	rabbitmq.PublishSimple("hello rabbit")
	fmt.Println("发送成功")
}
