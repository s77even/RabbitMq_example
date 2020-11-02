package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// url格式   “amqp：//账号：密码@地址：端口/vhost”
const MQURL = "amqp://imoocuser:imoocuser@127.0.0.1:5672/imooc"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机
	Exchange string
	//key
	Key string
	//连接信息
	Mqurl string
}

//NewRabbitMQ 创建RabbitMQ结构体实例
func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     MQURL,
	}
	var err error
	// 创建连接
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "创建连接错误……")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败……")
	return rabbitmq
}

//Destroy 断开连接
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

//failOnErr 错误处理
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s,", err, message)
		panic(fmt.Sprintf( "%s:%s,", err, message))
	}
}

//NewRabbitMQSimple  Simple模式
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	//simple模式下 exchange key为空 只需要传入一个name
	return NewRabbitMQ(queueName, "", "")
}

// PublishSimple 简单模式下的生产端
func (r *RabbitMQ) PublishSimple(message string) {
	// 1 申请队列 不存在则创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, //持久化 重启后丢失
		false, //自动删除 没有消费者后是否自动删除队列中的消息
		false, //排他性 是否只有自己可见
		false, //阻塞 是否等待响应
		nil,   //额外属性
	)
	if err != nil {
		fmt.Println(err)
	}
	//2 发送消息到队列
	err = r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false, // 无符合消息队列是否把消息返回给发送者
		false, // 没有绑定消费者是否返回给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		fmt.Println(err)
	}
}

//
func (r *RabbitMQ) ConsumeSimple() {
	// 申请队列
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, //持久化 重启后丢失
		false, //自动删除 没有消费者后是否自动删除队列中的消息
		false, //排他性 是否只有自己可见
		false, //阻塞 是否等待响应
		nil,   //额外属性
	)
	if err != nil {
		fmt.Println(err)
	}
	// 接收消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		"",   // 用来区分多个消费者
		true, // 是否自动应答ack
		false,
		false, //不能将同一个conn中的发送的消息传递给这个conn中的消费者
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	//启用协程处理消息
	go func() {
		for m := range msgs {
			// 消息处理逻辑
			log.Printf("received a message : %s", m.Body)
		}
	}()
	log.Printf("[*]waiting for messages...")
	<-forever // 阻塞
}

//NewRabbitMQPubSub 闯将订阅模式下的rabbit实例
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, "")
	var err error
	// 获取conn
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "创建连接错误……")
	// 获取 cahnnel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败……")
	return rabbitmq
}

//PublishPub 订阅模式下的发送消息
func (r *RabbitMQ) PublishPub(message string) {
	//创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false, // 此exchange能否被客户端用来推送消息
		false,
		nil)
	r.failOnErr(err, "falied to declear an exchange")
	// 发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

//ConsumeSub 订阅模式下的消费消息
func (r *RabbitMQ) ConsumeSub() {
	// 创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false, // 此exchange能否被客户端用来推送消息
		false,
		nil)
	r.failOnErr(err, "falied to declear an exchange")
	// 创建队列
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "falied to declear an queue")
	// 绑定队列到exchange
	err = r.channel.QueueBind(
		q.Name,
		"",
		r.Exchange,
		false,
		nil)

	// 接收消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		"",   // 用来区分多个消费者
		true, // 是否自动应答ack
		false,
		false, //不能将同一个conn中的发送的消息传递给这个conn中的消费者
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	//启用协程处理消息
	go func() {
		for m := range msgs {
			// 消息处理逻辑
			log.Printf("received a message : %s", m.Body)
		}
	}()
	log.Printf("[*]waiting for messages...")
	<-forever // 阻塞
}

// Routing 模式
//NewRabbitMQRouting 创建Routing实例
func NewRabbitMQRouting(exchangeName, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	// 获取conn
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq...")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel...")
	return rabbitmq
}

//PublishRouting 发送消息
func (r *RabbitMQ) PublishRouting(message string) {
	//创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct", // direct !!
		true,
		false,
		false,
		false,
		nil)
	r.failOnErr(err, "failed to declare an exchange")

	//发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.Key, // routingkey
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

//ConsumeRouting 消费消息
func (r *RabbitMQ) ConsumeRouting() {
	// 创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an exchange")
	// 创建队列
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		false,
		false,
		nil)
	r.failOnErr(err, "failed to declare an queue")
	// 队列绑定到交换机
	err = r.channel.QueueBind(
		q.Name,
		r.Key, // 需要绑定key
		r.Exchange,
		false,
		nil)
	// 消费消息
	message, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)
	go func(){
		for m := range message{
			log.Printf("received a message %s:" ,m.Body)
		}
	}()
	<-forever
}
