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

func taxCPU(cpuID int, pctLoad int, duration time.Duration, finished chan bool) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	C.lock_os_thread(C.int(cpuID))
	if C.sched_getcpu() != C.int(cpuID) {
		fmt.Printf("failed to bind thread to cpu %d\n", cpuID)
		finished <- false
		return
	}
	unitMultiplier := 100
	globalStart := time.Now()
	for {
		cycleStart := time.Now()
		for {
			cycleElapsed := time.Now().Sub(cycleStart)
			if cycleElapsed >= time.Microsecond*time.Duration(pctLoad)*time.Duration(unitMultiplier) {
				break
			}
		}
		time.Sleep(time.Microsecond * (100 - time.Duration(pctLoad)) * time.Duration(unitMultiplier))
		globalElapsed := time.Now().Sub(globalStart)
		if globalElapsed >= duration {
			break
		}
	}
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
		go taxCPU(workerNumber, 30, 10*time.Second, finished)
	}
	fmt.Println("waiting for worker to finish")
	for i := 0; i < len(cpus); i++ {
		<-finished
	}
	fmt.Println("DONE WITH MAIN")
}
