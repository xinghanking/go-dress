package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"go-dress/config"
	"sync"
	"time"
)

var once sync.Once
var Client *redis.Client

func Init() *redis.Client {
	once.Do(func() {
		cfg, err := config.LoadConfig()
		if err != nil {
			panic(err)
		}
		Client = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.Db,
		})
		_, err = Client.Ping().Result()
		if err != nil {
			panic(err)
		}
	})
	return Client
}

func Set(key string, value interface{}, expire uint) {
	expiration := time.Duration(expire) * time.Second
	err := Client.Set(key, value, expiration).Err()
	if err != nil {
		panic(err)
	}
}

func Get(key string) interface{} {
	result, err := Client.Get(key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		panic(err)
	}
	return result
}

func Del(key string) error {
	return Client.Del(key).Err()
}
