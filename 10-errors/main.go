package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch := make(chan bool, 1)
	defer close(ch)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(3 * time.Second)
		fmt.Println("Отдельная горутина отвисла")
		ch <- false
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Произошел тик тикера")
			ch <- true
		case value := <-ch:
			fmt.Printf("Получено значение %t\n", value)
			wg.Wait()
			return
		}
	}
}
