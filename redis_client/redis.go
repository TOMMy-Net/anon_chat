package redis_client

import(
	"log"
	"fmt"
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