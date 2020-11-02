package main

import "imooc_RabbitMQ/RabbitMQ"

func main() {
	rabbit :=  RabbitMQ.NewRabbitMQSimple("imoocSimple")
	rabbit.ConsumeSimple()
}