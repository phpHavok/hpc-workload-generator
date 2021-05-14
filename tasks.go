package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
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

type memoryTask struct {
	task
	numBytes uint
	duration time.Duration
}

func (t cpuLoadTask) execute(status chan int) {
	if err := lockOSThread(t.cpuID); err != nil {
		log.Error(err)
	}
	unitMultiplier := 1000
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

func (t memoryTask) execute(status chan int) {
	ptr, err := allocateMemory(t.numBytes)
	if err != nil {
		log.Error(err)
	}
	time.Sleep(t.duration)
	if err == nil {
		releaseMemory(ptr)
	}
	status <- 0
}

func (t memoryTask) String() string {
	return fmt.Sprintf("memory: %v", t.args)
}
