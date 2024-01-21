package storage

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type Storage struct {
	Client *redis.Client
}

func NewStorage(opts *StorageOptions) *Storage {
	client := redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Panicf("Failed to connect to Redis: %v", err)
	}

	return &Storage{
		Client: client,
	}
}
