package worker

import (
	"fmt"
	"sync"
)

var (
	stopChan = make(chan struct{})
	wg       sync.WaitGroup
)

func Run() {
	//wg.Add(1)
	//defer wg.Done()
	for {
		select {
		case <-stopChan:
			return
		default:
			// run all worker
			r, err := ParseBomFile()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(r)
		}
	}

}

func Stop() {
	// stop all worker
	close(stopChan)
	fmt.Println("Stopping worker...")
}

func Wait() {
	wg.Wait()
	fmt.Println("Worker stopped.")
}
