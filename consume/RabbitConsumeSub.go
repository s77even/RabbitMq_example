package main

import "imooc_RabbitMQ/RabbitMQ"

func main() {
	rabbit :=  RabbitMQ.NewRabbitMQPubSub("newProduct")
	rabbit.ConsumeSub()
}