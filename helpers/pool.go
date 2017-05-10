// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE IS SAFE TO EDIT BY HAND BUT IDEALLY YOU SHOULDN'T

package helpers

import (
	"sync"
)

// NewWorkerPool This function starts off n worker goroutines and allows you to send work to them.
//    In order to close down the work pool, just close the chan that is returned.
//    In order to ensure all workers have finished, call Wait() on the returned WaitGroup.
func NewWorkerPool(n int) (chan<- func(), *sync.WaitGroup) {
	work := make(chan func(), n)
	var wait sync.WaitGroup
	wait.Add(n)

	for ; n > 0; n-- {
		go func() {
			for x := range work {
				x()
			}
			wait.Done()
		}()
	}

	return work, &wait
}
