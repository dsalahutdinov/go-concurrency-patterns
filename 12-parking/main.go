package main

import (
	"fmt"
	"sync"
	"time"
)

type ParkingLot struct {
	slots chan struct{}
}

func (p *ParkingLot) Park(carID int64) {
	fmt.Printf("Car %d is being parking\n", carID)
	p.slots <- struct{}{}
	fmt.Printf("Car %d was parked successfully\n", carID)

	//go func() {
	time.Sleep(time.Second)
	fmt.Printf("Car %d left parking\n", carID)
	<-p.slots
	//}()
}

func main() {
	parking := &ParkingLot{slots: make(chan struct{}, 3)}

	var wg sync.WaitGroup
	carIDs := []int64{1, 2, 3, 4, 5, 6}
	for _, carID := range carIDs {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			parking.Park(id)
		}(carID)
	}
	wg.Wait()

}
