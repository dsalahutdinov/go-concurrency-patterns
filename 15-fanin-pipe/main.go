package main

import "fmt"

type Transaction struct {
	ID     int64
	Amount float64
}

func pipeline1(source <-chan Transaction, doSomething func(i Transaction) Transaction) <-chan Transaction {
	out := make(chan Transaction)

	go func() {
		for i := range source {
			out <- doSomething(i)
		}
		defer close(out)
	}()
	return out
}

func main() {

	source := make(chan Transaction)
	go func() {
		data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		for _, x := range data {
			source <- Transaction{ID: int64(x), Amount: float64(x) + 0.1}
		}
		close(source)
	}()

	resultChannel := pipeline1(source, func(i Transaction) Transaction { return Transaction{ID: i.ID, Amount: i.Amount + 1} })
	for r := range resultChannel {
		fmt.Println(r)
	}
}
