package provider

import (
	"github.com/jmoiron/sqlx"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/database"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/repository"
)

var (
	db         *sqlx.DB
	dataStore  repository.DataStore
	kafkaAdmin *database.KafkaAdmin
)

func BootstrapGlobal(cfg *config.Config) {
	db = database.NewPostgres((*database.PostgresOptions)(cfg.Postgres))
	dataStore = repository.NewDataStore(db)
	kafkaAdmin = database.NewKafkaAdmin(&database.KafkaAdminOptions{Brokers: cfg.Kafka.Brokers})
}
