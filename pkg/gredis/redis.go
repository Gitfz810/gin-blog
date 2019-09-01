package gredis

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"

	"gin-blog/pkg/setting"
)

var RedisConn *redis.Pool

// Setup init the redis instance
func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:       setting.RedisSetting.MaxIdle,     // 最大空闲连接数
		MaxActive:     setting.RedisSetting.MaxActive,   // 给定时间内，允许分配的最大连接数
		IdleTimeout:   setting.RedisSetting.IdleTimeout, // 在给定时间内将会保持空闲状态，若达到时间限制则关闭连接
		// 提供创建和配置应用程序连接的一个函数
		Dial: func() (conn redis.Conn, e error) {
			conn, e = redis.Dial("tcp", setting.RedisSetting.Host)
			if e != nil {
				return nil, e
			}
			if setting.RedisSetting.Password != "" {
				if _, e = conn.Do("AUTH", setting.RedisSetting.Password); e != nil {
					conn.Close()
					return nil, e
				}
			}
			return
		},
		// 可选的应用程序检查健康功能
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

/*
封装一些基础指令功能
*/
// setex key second value
func Setex(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	_, err = conn.Do("SETEX", key, time, value)
	if err != nil {
		return err
	}
	return nil
}

// exists key
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

// get key
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// delete key
func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// batch delete func
func BatchDelete(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	pattern := strings.Join([]string{"*", key, "*"}, "")
	keys, err := redis.Strings(conn.Do("KEYS", pattern))
	if err != nil {
		return err
	}
	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}