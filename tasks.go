package main

import (
	"fmt"
	"time"
)

type executable interface {
	execute(status chan int)
}

type task struct {
	args []string
}

type cpuLoadTask struct {
	task
	cpuID    int
	pctLoad  int
	duration time.Duration
}

func (t cpuLoadTask) execute(status chan int) {
	if err := lockOSThread(t.cpuID); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	unitMultiplier := 100
	globalStart := time.Now()
	for {
		cycleStart := time.Now()
		for {
			cycleElapsed := time.Now().Sub(cycleStart)
			if cycleElapsed >= time.Microsecond*time.Duration(t.pctLoad)*time.Duration(unitMultiplier) {
				break
			}
		}
		time.Sleep(time.Microsecond * (100 - time.Duration(t.pctLoad)) * time.Duration(unitMultiplier))
		globalElapsed := time.Now().Sub(globalStart)
		if globalElapsed >= t.duration {
			break
		}
	}
	status <- 0
}

func (t cpuLoadTask) String() string {
	return fmt.Sprintf("cpuload: %v", t.args)
}
