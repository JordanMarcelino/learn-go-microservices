package database

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisOptions struct {
	Addrs           []string
	Password        string
	DialTimeout     int
	ReadTimeout     int
	WriteTimeout    int
	MaxIdleConn     int
	MaxActiveConn   int
	MaxConnLifetime int
}

func InitRedisCluster(opt *RedisOptions) *redis.ClusterClient {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           opt.Addrs,
		Password:        opt.Password,
		DialTimeout:     time.Duration(opt.DialTimeout) * time.Second,
		ReadTimeout:     time.Duration(opt.ReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(opt.WriteTimeout) * time.Second,
		MinIdleConns:    opt.MaxIdleConn,
		ConnMaxLifetime: time.Duration(opt.MaxConnLifetime) * time.Minute,
		MaxActiveConns:  opt.MaxActiveConn,
	})

	ctx := context.Background()
	status, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to redis cluster: %v", err)
	}
	log.Println("connected to redis cluster:", status)

	return rdb
}
