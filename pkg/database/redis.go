package database

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClusterOptions struct {
	Addrs           []string
	Password        string
	DialTimeout     int
	ReadTimeout     int
	WriteTimeout    int
	MinIdleConn     int
	MaxIdleConn     int
	MaxActiveConn   int
	MaxConnLifetime int
}

type RedisOptions struct {
	Addr            string
	DialTimeout     int
	ReadTimeout     int
	WriteTimeout    int
	MinIdleConn     int
	MaxIdleConn     int
	MaxActiveConn   int
	MaxConnLifetime int
}

func NewRedisCluster(opt *RedisClusterOptions) *redis.ClusterClient {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           opt.Addrs,
		Password:        opt.Password,
		DialTimeout:     time.Duration(opt.DialTimeout) * time.Second,
		ReadTimeout:     time.Duration(opt.ReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(opt.WriteTimeout) * time.Second,
		MinIdleConns:    opt.MinIdleConn,
		MaxIdleConns:    opt.MaxIdleConn,
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

func NewRedis(opt *RedisOptions) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:            opt.Addr,
		DialTimeout:     time.Duration(opt.DialTimeout) * time.Second,
		ReadTimeout:     time.Duration(opt.ReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(opt.WriteTimeout) * time.Second,
		MinIdleConns:    opt.MinIdleConn,
		MaxIdleConns:    opt.MaxIdleConn,
		ConnMaxLifetime: time.Duration(opt.MaxConnLifetime) * time.Minute,
		MaxActiveConns:  opt.MaxActiveConn,
	})

	ctx := context.Background()
	status, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	log.Println("connected to redis:", status)

	return rdb
}
