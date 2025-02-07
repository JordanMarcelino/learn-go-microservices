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
	KafkaRetryDelay = 3
	KafkaRetryLimit = 3
)
