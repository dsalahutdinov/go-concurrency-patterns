package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	// Create an errgroup with a context
	g, ctx := errgroup.WithContext(context.Background())

	// Add goroutines to the group
	g.Go(func() error {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("Task 1 completed")
			return nil
		case <-ctx.Done():
			fmt.Println("Task 1 canceled")
			return ctx.Err() // Return the context cancellation error
		}
	})

	g.Go(func() error {
		select {
		case <-time.After(500 * time.Millisecond):
			fmt.Println("Task 2 completed with an error")
			return fmt.Errorf("error from task 2")
		case <-ctx.Done():
			fmt.Println("Task 2 canceled")
			return ctx.Err()
		}
	})

	// Wait for all goroutines to complete
	if err := g.Wait(); err != nil {
		fmt.Printf("Group finished with error: %v\n", err)
	} else {
		fmt.Println("All tasks completed successfully")
	}
}
