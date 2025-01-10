package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
)

type User struct {
	name      string `redis:"name"`
	last_name string `redis:"last_name"`
	email     string `redis:"email"`
	age       string `redis:"age"`
}

func printMap(conn redis.Conn, hashKey string) {
	name, err := redis.String(conn.Do("HGET", hashKey, "email"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)
}
func printMap2(conn redis.Conn, hashKey string) {
	hashMap, err := redis.StringMap(conn.Do("HGETALL", hashKey))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v", hashMap)
}
func printMap3(conn redis.Conn, hashKey string) {
	values, err := redis.Values(conn.Do("HGETALL", hashKey))
	if err != nil {
		log.Fatal(err)
	}
	var user User
	err = redis.ScanStruct(values, &user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v", user)
}
func main() {
	ctx := context.Background()
	conn, err := redis.DialContext(ctx, "tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	hashKey := gofakeit.UUID()
	users := map[string]string{
		"name":      gofakeit.Name(),
		"last_name": gofakeit.LastName(),
		"age":       strconv.Itoa(gofakeit.IntRange(0, 100)),
		"email":     gofakeit.Email(),
	}
	for key, value := range users {
		_, err := conn.Do("HSET", hashKey, key, value)
		if err != nil {
			log.Fatal(err)
		}
	}
	//printMap(conn, hashKey)
	//fmt.Printf("\n")
	//printMap2(conn, hashKey)
	//fmt.Printf("\n")
	printMap3(conn, hashKey)
}
