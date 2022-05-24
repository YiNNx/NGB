package model

import (
	"github.com/go-redis/redis/v7"
	"ngb/util/log"
	"time"
)

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(
		&redis.Options{
			Addr: "redis:6379",
		},
	)

	err := redisClient.Ping().Err()
	if err != nil {
		log.Logger.Error(err)
	} else {
		log.Logger.Printf("Redis server connected")
	}
}

func redisErrHandler(err error) error {
	log.Logger.Printf("Redis error: %v", err)
	return err
}

func RedisSet(key string, value interface{}, exp time.Duration) error {
	log.Logger.Printf(key, ":", value, " ; ", exp)
	err := redisClient.Set(key, value, exp).Err()
	if err != nil {
		return redisErrHandler(err)
	}
	return nil
}

func RedisPush(key string, list []string) error {
	err := redisClient.LPush(key, list).Err()
	if err != nil {
		return redisErrHandler(err)
	}
	return nil
}

func RedisGet(key string) (value interface{}, err error) {

	value, err = redisClient.Get(key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, redisErrHandler(err)
	}
	log.Logger.Printf(key, ":", value.(string))
	return value, nil
}

func RedisDelete(key string) error {
	err := redisClient.Del(key).Err()
	if err != nil && err != redis.Nil { // Will it?
		return redisErrHandler(err)
	}
	return nil
}

func RedisSetFailList(failList []string, subject string, text string) {
	err := RedisPush("fail-list", failList)
	if err != nil {
		log.Logger.Error(err)
	}
	err = RedisSet("subject", subject, 0)
	if err != nil {
		log.Logger.Error(err)
	}
	err = RedisSet("text", text, 0)
	if err != nil {
		log.Logger.Error(err)
	}
}

func RedisReadFailList() ([]string, string, string) {
	list, err := redisClient.LRange("fail-list", 0, -1).Result()
	if err != nil {
		log.Logger.Error(err)
	}
	if len(list) == 0 {
		return nil, "", ""
	}
	log.Logger.Info(list)
	subject, err := RedisGet("subject")
	if err != nil {
		log.Logger.Error(err)
	}
	text, err := RedisGet("text")
	if err != nil {
		log.Logger.Error(err)
	}
	if err := RedisDelete("fail-list"); err != nil {
		log.Logger.Error(err)
	}
	if err := RedisDelete("subject"); err != nil {
		log.Logger.Error(err)
	}
	if err := RedisDelete("text"); err != nil {
		log.Logger.Error(err)
	}
	return list, subject.(string), text.(string)
}
