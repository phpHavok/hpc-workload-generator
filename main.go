package main

// #define _GNU_SOURCE
// #include <sched.h>
// #include "helper.h"
import "C"

import (
	"flag"
	"fmt"
	"runtime"

	log "github.com/sirupsen/logrus"
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
	scheduleFile := flag.String("schedule-file", "", "path to the schedule file to execute, or empty for stdin")
	cgroupsRootPath := flag.String("cgroups-root", "/sys/fs/cgroup", "path to the root of the cgroupsv1 hierarchy")
	flag.Parse()
	schedule, err := loadSchedule(*scheduleFile, *cgroupsRootPath)
	if err != nil {
		log.Fatalf("failed to load schedule: %v", err)
	}
	if err := schedule.execute(); err != nil {
		log.Fatalf("failed to execute schedule: %v", err)
	}
}
