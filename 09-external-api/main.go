// Как сделать влияние функции getDiscount более управляемым
package main

import (
	"context"
	"fmt"
	"time"
)

var defaultTimeout = 1 * time.Second

func getDiscount() float64 {
	time.Sleep(2000 * time.Millisecond)
	return 12.0
}

func getDiscountWithTimeout(ctx context.Context) (float64, error) {

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
		defer cancel()
	}

	res := make(chan float64)
	defer close(res)
	go func() {
		res <- getDiscount()
	}()
	select {
	case <-ctx.Done():
		return 0.0, ctx.Err()
	case val := <-res:
		return val, nil
	}
}

func main() {
	result, err := getDiscountWithTimeout(context.Background())
	fmt.Printf("Ваша скидка: %v, err: %s\n", result, err)
}
