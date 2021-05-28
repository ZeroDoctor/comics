package util

import (
	"fmt"
	"sync"
)

type Process struct {
	Args []interface{}
	Name string
	Fn   func(Process)
}

func WorkerPool(mp []Process, workers int) {
	var wg sync.WaitGroup
	currQueue := 0
	waitLimit := workers

	if len(mp) < workers {
		waitLimit = len(mp)
	}
	wg.Add(waitLimit)

	for i, proc := range mp {
		go func(proc Process) {
			fmt.Println("starting goprocess")
			proc.Fn(proc)
			wg.Done()
		}(proc)

		currQueue++
		if currQueue >= workers && i != len(mp) {
			currQueue = 0
			// ppt.Infof("currently: (%d/%d)\n", i+1, len(mp))

			waitLimit := workers
			if len(mp)-(i+1) < workers {
				waitLimit = (len(mp) - 1) - i
			}
			wg.Wait()
			wg.Add(waitLimit)
		}
	}

	// ppt.Infof("currently: (%d/%d)\n", i+1, len(mp))
	wg.Wait()
}
