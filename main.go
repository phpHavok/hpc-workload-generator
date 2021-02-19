package main

// #define _GNU_SOURCE
// #include <sched.h>
// #include "helper.h"
import "C"

import (
	"fmt"
	"runtime"
	"time"
)

func worker(id int, finished chan bool) {
	runtime.LockOSThread()
	C.lock_os_thread(C.int(2))
	globalStart := time.Now()
	for {
		cycleStart := time.Now()
		for {
			cycleElapsed := time.Now().Sub(cycleStart)
			if cycleElapsed >= time.Microsecond*3000 {
				break
			}
		}
		time.Sleep(time.Microsecond * 7000)
		globalElapsed := time.Now().Sub(globalStart)
		if globalElapsed >= time.Second*15 {
			break
		}
	}
	fmt.Println("Done")
	finished <- true
}

func main() {
	finished := make(chan bool)
	fmt.Println("HEllo world")
	go worker(1, finished)
	fmt.Println("waiting for worker to finish")
	<-finished
	fmt.Println("DONE WITH MAIN")
}
