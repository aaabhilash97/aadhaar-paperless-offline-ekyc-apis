package aadhaarcache

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(opt *redis.Options) RedisCache {
	return RedisCache{
		rdb: redis.NewClient(opt),
	}
}

func (c RedisCache) SaveSession(session string) (hash string, err error) {
	hash = fmt.Sprintf("%x", md5.Sum([]byte(session)))
	{
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		if _, err = c.rdb.HSet(ctx, hash, "session", session).Result(); err != nil {
			return
		}
	}
	{
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		if _, err = c.rdb.Expire(ctx, hash, time.Minute*10).Result(); err != nil {
			return
		}
	}
	return
}

func (c RedisCache) GetSession(hash string) (session string, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if session, err = c.rdb.HGet(ctx, hash, "session").Result(); err != nil {
		return
	}

	return
}

func (c RedisCache) SaveData(hash string, data interface{}) (err error) {
	dataJson, err := json.Marshal(data)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if _, err = c.rdb.HSet(ctx, hash, "data",
		base64.StdEncoding.EncodeToString(dataJson)).Result(); err != nil {
		return
	}

	return
}

func (c RedisCache) GetData(hash string, v interface{}) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	var jsonBase64 string
	if jsonBase64, err = c.rdb.HGet(ctx, hash, "data").Result(); err != nil {
		return
	}

	data, err := base64.StdEncoding.DecodeString(jsonBase64)
	if err != nil {
		return
	}
	if err = json.Unmarshal(data, v); err != nil {
		return
	}
	return

}

func (c RedisCache) IsNotFoundError(err error) bool {
	return err == redis.Nil
}
