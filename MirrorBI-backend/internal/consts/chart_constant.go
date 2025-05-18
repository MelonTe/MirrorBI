package consts

const (
	ChartStatusSucceed    = "succeed"
	ChartStatusFailed     = "failed"
	ChartStatusRunning    = "running"
	ChartStatusWait       = "wait"
	ChartStatusNotSucceed = "not_succeed"
	ChartConsumerName     = "ChartService" // 消费者名称

	MQExchangeName = "mrbi" // 交换机名称
	MQQueueName    = "mrbi" // 队列名称
	MQRoutingKey   = "mrbi" // 路由键名称
)
