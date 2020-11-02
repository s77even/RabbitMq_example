package main

import "imooc_RabbitMQ/RabbitMQ"

func main() {
	imooc_one :=  RabbitMQ.NewRabbitMQRouting("exImooc","imooc_one")
	imooc_one.ConsumeRouting()
}