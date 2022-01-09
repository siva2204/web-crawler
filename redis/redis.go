package redis_crawler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var (
	ErrInvalidKey = errors.New("key doesn't exist")
	ErrEmptyValue = errors.New("value is empty")
	ErrRedis      = errors.New("redis error")
	ErrFormat     = errors.New("json marshall error")
)

type RedisClient struct {
	RDB *redis.Client
	ctx context.Context
}

type redisPayLoad struct {
	Payload []string `json:"Payload"`
}

var Client *RedisClient

func CreateClient(host string, port string) {
	ctx := context.Background()
	addr := fmt.Sprintf("%s:%s", host, port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	// testing connection with redis server
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(fmt.Errorf("error connecting with redis: %+v", err))
	}

	Client = &RedisClient{RDB: rdb, ctx: ctx}
}

func (client *RedisClient) Insert(key string, value []string) error {

	payload := &redisPayLoad{Payload: value}

	val, err := json.Marshal(payload)

	if err != nil {
		return ErrFormat
	}

	if err = client.RDB.Set(client.ctx, key, val, 0).Err(); err != nil {
		return ErrRedis
	}

	return nil
}

func (client *RedisClient) Append(key string, value []string) error {
	val, err := client.RDB.Get(client.ctx, key).Result()
	switch {
	case err == redis.Nil:
		return ErrInvalidKey
	case err != nil:
		return ErrRedis
	case val == "":
		// val = '[]'
	}

	fmt.Println(val)

	x := []byte(val)
	var data redisPayLoad

	err = json.Unmarshal(x, &data)

	if err != nil {
		return ErrFormat
	}
	fmt.Println(data.Payload)

	newPayload := append(data.Payload, value...)

	return client.Insert(key, newPayload)
}
