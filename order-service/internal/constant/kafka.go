package constant

const (
	ProductCreatedTopic = "product-created"
	ProductUpdatedTopic = "product-updated"
	ProductDeletedTopic = "product-deleted"

	OrderCreatedTopic                  = "order-created"
	OrderCreatedTopicNumPartitions     = 3
	OrderCreatedTopicReplicationFactor = 1

	OrderCancelledTopic                  = "order-cancelled"
	OrderCancelledTopicNumPartitions     = 3
	OrderCancelledTopicReplicationFactor = 1
)

const (
	ProductCreatedConsumerGroup = "product-created-consumer-group"

	ProductUpdatedConsumerGroup = "product-updated-consumer-group"

	ProductDeletedConsumerGroup = "product-deleted-consumer-group"
)

const (
	KafkaProducerRetryDelay = 3
	KafkaProducerRetryLimit = 3

	KafkaConsumerRetryDelay = 2
	KafkaConsumerRetryLimit = 3
)
