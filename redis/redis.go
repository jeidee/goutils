package redis

import (
	"fmt"

	redis "gopkg.in/redis.v5"
)

// New 함수는 새로운 redis client 객체를 반환한다.
func New(host string, port int, password string, defaultDB int) *redis.Client {
	addr := fmt.Sprintf("%s:%d", host, port)
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       defaultDB,
	})
}

// Ping 함수는 redis의 연결 상태를 반환한다.
func Ping(client *redis.Client) error {
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	if err != nil {
		return err
	}
	return nil
}

func Get(client *redis.Client, key string) {

	redis.NewStringCmd("GET", key)
}
