package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
	"task/floodControl"
	"time"
)

func Sandbox(fc FloodControl, numWorkers int) {

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			userID := int64(workerID % 3) // Пример: разные пользователи для разных горутин

			for j := 0; j < 7; j++ {
				passed, err := fc.Check(context.Background(), userID)
				if err != nil {
					fmt.Println(err)
					return
				}
				if passed {
					fmt.Printf("Горутина %d, пользователь %d, Check вызов %d: Проверка пройдена\n", workerID, userID, j+1)
				} else {
					fmt.Printf("Горутина %d, пользователь %d, Check вызов %d: Превышен лимит запросов\n", workerID, userID, j+1)
				}
			}
		}(i)
	}
	wg.Wait()
}

func main() {
	//В случае использования redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	fc := floodControl.NewFloodControlRedisImpl(client, 5, 10*time.Second, "flood_control")

	//в случае использования решения на мапе
	//fc := floodControl.NewFloodControlMapImpl(10, 5)

	Sandbox(fc, 7)
}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}
