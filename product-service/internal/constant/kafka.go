package constant

const (
	ProductCreatedTopic                  = "product-created"
	ProductCreatedTopicNumPartitions     = 3
	ProductCreatedTopicReplicationFactor = 1

	ProductUpdatedTopic                  = "product-updated"
	ProductUpdatedTopicNumPartitions     = 3
	ProductUpdatedTopicReplicationFactor = 1

	ProductDeletedTopic                  = "product-deleted"
	ProductDeletedTopicNumPartitions     = 3
	ProductDeletedTopicReplicationFactor = 1
)

const (
	OrderCreatedTopic   = "order-created"
	OrderCancelledTopic = "order-cancelled"
)

const (
	OrderCreatedConsumerGroup   = "order-created-consumer-group"
	OrderCancelledConsumerGroup = "order-cancelled-consumer-group"
)

const (
	KafkaProducerRetryDelay = 3
	KafkaProducerRetryLimit = 3

	KafkaConsumerRetryDelay = 2
	KafkaConsumerRetryLimit = 3
)
