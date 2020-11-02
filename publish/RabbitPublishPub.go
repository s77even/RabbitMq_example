package main

import (
	"fmt"
	"imooc_RabbitMQ/RabbitMQ"
	"strconv"
	"time"
)

func main(){
	rabbitmq:= RabbitMQ.NewRabbitMQPubSub("newProduct")
	for i:=0 ; i<20 ; i++{
		rabbitmq.PublishPub("订阅模式第"+ strconv.Itoa(i)+"条消息")
		fmt.Println("订阅模式第"+ strconv.Itoa(i)+"条消息")
		time.Sleep(time.Second)
	}
}
