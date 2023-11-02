package redis_client

import (
	"fmt"
	"log"
	redis "github.com/go-redis/redis"
)
type Redis struct{
	Client *redis.Client
}
func Create_client() Redis{
	// Создание клиента Redis
	client := Redis{Client: redis.NewClient(&redis.Options{
	 Addr:     "127.0.0.1:6379", // адрес и порт Redis-сервера
	 Password: "",               // пароль Redis-сервера, если есть
	 DB:       0,                // номер базы данных Redis
	})}
   
	// Проверка соединения с Redis
	_, err := client.Client.Ping().Result()
	if err != nil {
	 log.Fatal(err)
	}
	fmt.Println("Соединение с Redis установлено")
	return client
}

func (r Redis)Queue_add(data int64){
	err := r.Client.RPush("search", data).Err()
	if err != nil {
	 log.Fatal(err)
	}
}

func (r Redis)Queue_pop() string {
	element, err := r.Client.LPop("search").Result()
	if err == redis.Nil{
		return ""
	}else if err != nil {
		panic(err)
	}
	return element
}

func (r Redis)Queue_rem(data int64)  {
	r.Client.LRem("search", 0, data)
}