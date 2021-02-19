package main

// #define _GNU_SOURCE
// #include <sched.h>
// #include "helper.h"
import "C"

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/phpHavok/hpc-workload-generator/cgroups"
)

func worker(id int, finished chan bool) {
	runtime.LockOSThread()
	C.lock_os_thread(C.int(id))
	fmt.Println("Running process: ", id, " on CPU: ", C.sched_getcpu())
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
	fmt.Println("Done with process: ", id, " on CPU: ", C.sched_getcpu())
	finished <- true
}

func main() {
	processCgroups, err := cgroups.LoadProcessCgroups(os.Getpid())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	cpus, err := processCgroups.Cpuset.GetCpus()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	fmt.Println("CPUSET: ", cpus)

	runtime.GOMAXPROCS(len(cpus))
	finished := make(chan bool, len(cpus))
	for _, workerNumber := range cpus {
		go worker(workerNumber, finished)
	}
	fmt.Println("waiting for worker to finish")
	for i := 0; i < len(cpus); i++ {
		<-finished
	}
	fmt.Println("DONE WITH MAIN")
}
