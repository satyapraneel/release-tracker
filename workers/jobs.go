package workers

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func Wait() {
	wg.Wait()
}

type Job struct {
	Action  func(map[string]interface{})
	Payload map[string]interface{}
}

func (job Job) Dispatch() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		job.Action(job.Payload)

	}()
}

func PrintPayload(payload map[string]interface{}) {
	time.Sleep(15)
	fmt.Println(payload)
}
