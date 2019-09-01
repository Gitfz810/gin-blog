package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func main() {
	server := "127.0.0.1:6379"

	conn, err := redis.Dial("tcp", server)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	_, err = conn.Do("SET", "aaa", "test redigo!")
	if err != nil {
		fmt.Println(err)
		return
	}

	exists, err := redis.Bool(conn.Do("EXISTS", "aaa"))
	if err != nil {
		fmt.Println(false)
		return
	}
	fmt.Println(exists)
}
