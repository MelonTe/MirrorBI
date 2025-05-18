package mq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"mrbi/config"
	"mrbi/internal/consts"
)

var connPool *ChannelPool

// channel池
type ChannelPool struct {
	conn *amqp.Connection
	pool chan *amqp.Channel
}

// 初始化，创建交换机、队列、初始化连接池
func init() {
	//初始化连接池
	cfg := config.LoadConfig()
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.RabbitMQ.UserName, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host, cfg.RabbitMQ.Port))
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ")
	}
	connPool = &ChannelPool{
		conn: conn,
		pool: make(chan *amqp.Channel, 5), // 设置池的大小为5
	}
	//创建交换机和队列
	ch, err := conn.Channel()
	if err != nil {
		log.Panic("Failed to open a channel")
	}
	err = ch.ExchangeDeclare(
		consts.MQExchangeName, // 交换机名称
		"direct",              // 交换机类型
		true,                  // 是否持久化
		false,                 // 是否自动删除
		false,                 // 是否内部使用
		false,                 // 是否等待服务器确认
		nil,                   // 额外参数
	)
	failOnError(err, "Failed to declare an exchange")
	//声明队列
	_, err = ch.QueueDeclare(
		consts.MQQueueName, // 队列名称
		true,               // 是否持久化
		false,              // 是否自动删除
		false,              // 是否排他
		false,              // 是否等待服务器确认
		nil,                // 额外参数
	)
	failOnError(err, "Failed to declare a queue")
	//绑定交换机和队列
	err = ch.QueueBind(
		consts.MQQueueName,    // 队列名称
		consts.MQRoutingKey,   // 路由键名称
		consts.MQExchangeName, // 交换机名称
		false,                 // 是否等待服务器确认
		nil,                   // 额外参数
	)
	failOnError(err, "Failed to bind queue to exchange")
	//将通道放入连接池
	for i := 0; i < cap(connPool.pool); i++ {
		channel, err := conn.Channel()
		if err != nil {
			log.Panic("Failed to open a channel")
		}
		err = channel.Qos(20, 0, false)
		failOnError(err, "Failed to set QoS")
		connPool.pool <- channel
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func GetChannelPool() *ChannelPool {
	return connPool
}
func GetChannel() *amqp.Channel {
	ch := <-connPool.pool
	return ch
}
func ReleaseChannel(ch *amqp.Channel) {
	//将通道放回连接池
	connPool.pool <- ch
}

func (connPool *ChannelPool) PublishMessage(message []byte) error {
	//从连接获取一个通道
	ch := <-connPool.pool
	//释放连接
	defer func() {
		connPool.pool <- ch
	}()
	//发布消息
	err := ch.Publish(
		consts.MQExchangeName, // 交换机名称
		consts.MQRoutingKey,   // 路由键名称
		false,                 // 是否等待服务器确认
		false,                 // 是否强制发布
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	return err
}
