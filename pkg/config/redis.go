package config

import "github.com/go-redis/redis/v8"

type Redis struct {
	HOST     string
	PASSWORD string
	PORT     string
	DB       int
}

func NewRedis(data Redis) *redis.Client{
	options := &redis.Options{
		Addr: data.HOST + ":" + data.PORT,
		DB: data.DB,
		Password: data.PASSWORD,
	}

	client := redis.NewClient(options)
	// defer client.Close()

	return client
}