package aadhaarcache

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"time"

	"github.com/patrickmn/go-cache"
)

type MemCache struct {
	db *cache.Cache
}

func NewMemCache() *MemCache {
	return &MemCache{
		db: cache.New(10*time.Minute, 10*time.Minute),
	}
}

func (c MemCache) SaveSession(session string) (hash string, err error) {
	hash = fmt.Sprintf("%x", md5.Sum([]byte(session)))
	c.db.Set(fmt.Sprintf("session:%s", hash), session, 10*time.Minute)
	return
}

var memNilError error = fmt.Errorf("Not found")

func (c MemCache) GetSession(hash string) (session string, err error) {

	if val, ok := c.db.Get(fmt.Sprintf("session:%s", hash)); ok {
		if session, ok := val.(string); ok {
			return session, nil
		}
	}
	return session, memNilError
}

func (c MemCache) SaveData(hash string, data interface{}) (err error) {
	c.db.Set(fmt.Sprintf("data:%s", hash), data, 10*time.Minute)
	return
}

func (c MemCache) GetData(hash string, v interface{}) (err error) {

	val, ok := c.db.Get(fmt.Sprintf("data:%s", hash))
	if ok {
		v := reflect.ValueOf(v).Elem()
		v.Set(reflect.ValueOf(val))
		return
	}
	return memNilError

}

func (c MemCache) IsNotFoundError(err error) bool {
	return err == memNilError
}
