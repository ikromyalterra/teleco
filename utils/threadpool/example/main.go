package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sepulsa/teleco/utils/threadpool"
)

func main() {

	for i := 1; i <= 5; i++ {
		task := &myTask{ID: i}
		threadpool.Run("test", 1, 2, task)
	}

}

type myTask struct {
	ID int
}

func (m *myTask) Run() {
	fmt.Println("Running my task ", m.ID)
	sleep := rand.Intn(6)
	time.Sleep(time.Duration(sleep) * time.Second)
}

func (m *myTask) RunWhenTimeout() {
	fmt.Println("Running when Timeout my task ", m.ID)
}

func (m *myTask) RunAfterTimeout() {
	fmt.Println("Running after Timeout my task ", m.ID)
}

func (m *myTask) RunWhenFull() {
	fmt.Println("Running when Full my task ", m.ID)
}
