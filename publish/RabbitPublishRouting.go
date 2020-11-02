package main

import (
	"fmt"
	"imooc_RabbitMQ/RabbitMQ"
	"strconv"
	"time"
)

func main(){
	imoocOne := RabbitMQ.NewRabbitMQRouting("exImooc","imooc_one")
	imoocTwo := RabbitMQ.NewRabbitMQRouting("exImooc","imooc_two")

	for i:=0 ; i <=10 ; i++{
		imoocOne.PublishRouting("hello imooc One"+strconv.Itoa(i))
		imoocTwo.PublishRouting("hello imooc Two"+strconv.Itoa(i))
		fmt.Println(i)
		time.Sleep(time.Second)
	}
}
