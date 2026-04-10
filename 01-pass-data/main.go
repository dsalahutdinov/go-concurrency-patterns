package main

import "fmt"

func main(){
	ch := make(chan int)

	go func() {
		ch <- 43
	}()

	value := <- ch

	fmt.Println(value)
}