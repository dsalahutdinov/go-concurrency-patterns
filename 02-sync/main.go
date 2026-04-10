package main

import (
	"fmt"
	"time"
)

func worker(done chan bool) {
	fmt.Println("Worker started")
	time.Sleep(time.Second)
	fmt.Println("Worker finishied")
	done <- true
}

func main() {

	doneChannel := make(chan bool, 5)

	for i := 0; i < 5; i++ {
		go worker(doneChannel)
	}

	counter := 0
	for {

		<-doneChannel
		counter += 1
		fmt.Printf("Counter %d\n", counter)
		if counter == 5 {
			close(doneChannel)
			fmt.Println("Selected")
			return
		}
	}
}
