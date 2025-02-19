package provider

import (
	"github.com/bsm/redislock"
	"github.com/jmoiron/sqlx"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/repository"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/database"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

var (
	db         *sqlx.DB
	rdb        *redis.ClusterClient
	rdl        *redislock.Client
	rabbitmq   *amqp.Connection
	dataStore  repository.DataStore
	kafkaAdmin *database.KafkaAdmin
)

func BootstrapGlobal(cfg *config.Config) {
	db = database.NewPostgres((*database.PostgresOptions)(cfg.Postgres))
	rdb = database.NewRedisCluster((*database.RedisClusterOptions)(cfg.Redis))
	rdl = redislock.New(rdb)
	rabbitmq = database.NewAMQP((*database.AmqpOptions)(cfg.Amqp))
	dataStore = repository.NewDataStore(db)
	kafkaAdmin = database.NewKafkaAdmin(&database.KafkaAdminOptions{Brokers: cfg.Kafka.Brokers})
}
