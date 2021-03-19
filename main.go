package main

// #define _GNU_SOURCE
// #include <sched.h>
// #include "helper.h"
import "C"

import (
	"fmt"
	"os"
	"runtime"
)

func lockOSThread(cpuID int) error {
	runtime.LockOSThread()
	C.lock_os_thread(C.int(cpuID))
	if C.sched_getcpu() != C.int(cpuID) {
		return fmt.Errorf("failed to bind thread to cpu %d", cpuID)
	}
	return nil
}

func main() {
	schedule, err := loadSchedule("timeline.in")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	if err := schedule.execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	fmt.Println("DONE WITH MAIN")
}
