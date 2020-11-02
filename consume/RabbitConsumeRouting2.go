package main

import "imooc_RabbitMQ/RabbitMQ"

func main() {
	imooc_two :=  RabbitMQ.NewRabbitMQRouting("exImooc","imooc_two")
	imooc_two.ConsumeRouting()
}