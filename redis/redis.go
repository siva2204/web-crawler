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

func CreateClient(host string, port int) {
	ctx := context.Background()
	addr := fmt.Sprintf("%s:%d", host, port)
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
		fmt.Println("err", err)
		return ErrRedis
	}

	return nil
}

// creates a set for a key and adds the value in set
func (client *RedisClient) Append(key string, value string) error {

	if err := client.RDB.SAdd(client.ctx, key, value).Err(); err != nil {
		fmt.Errorf("Error adding url %s to key %s", value, key)
		return err
	}

	return nil
}

func (client *RedisClient) Get(key string) ([]string, error) {
	val, err := client.RDB.Get(client.ctx, key).Result()
	switch {
	case err == redis.Nil:
		return []string{}, ErrInvalidKey
	case err != nil:
		fmt.Println("err", err)
		return []string{}, ErrRedis
	case val == "":
		return []string{}, ErrEmptyValue
	}

	x := []byte(val)
	var data redisPayLoad

	err = json.Unmarshal(x, &data)

	if err != nil {
		return []string{}, ErrFormat
	}

	return data.Payload, nil

}

func (client *RedisClient) AutoComplete(key string) ([]string, error) {
	var cursor uint64
	fullkeys := []string{}
	for {
		var keys []string
		var err error
		keys, cursor, err = client.RDB.Scan(client.ctx, cursor, key+"*", 10).Result()
		if err != nil {
			fmt.Println("err", err)
			return []string{}, ErrRedis
		}
		if cursor == 0 {
			break
		}
		fullkeys = append(fullkeys, keys...)
	}
	return fullkeys, nil
}
func (client *RedisClient) GetAll() ([]string, error) {

	dataSet := []string{}

	iter := client.RDB.Scan(client.ctx, 0, "*", 0).Iterator()
	for iter.Next(client.ctx) {
		dataSet = append(dataSet, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return dataSet, err
	}

	return dataSet, nil
}

func (client *RedisClient) GetMany(keys []string) (map[string][]string, error) {
	var hm map[string][]string = make(map[string][]string)

	for _, x := range keys {
		val, err := client.GetSetValues(x)

		if err != nil {
			fmt.Errorf("error getting data form redis for key %s", x)
		}
		hm[x] = val
	}
	return hm, nil
}

// get all the values in the set for a given key
func (client *RedisClient) GetSetValues(key string) ([]string, error) {
	val, err := client.RDB.SMembers(client.ctx, key).Result()

	switch {
	case err == redis.Nil:
		return []string{}, ErrInvalidKey
	case err != nil:
		return []string{}, ErrRedis
	}

	return val, nil
}
