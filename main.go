package main

// #define _GNU_SOURCE
// #include <sched.h>
// #include "helper.h"
import "C"

import (
	"flag"
	"fmt"
	"runtime"
	"unsafe"

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

func allocateMemory(numBytes uint) (unsafe.Pointer, error) {
	ptr := C.allocate_memory(C.ulong(numBytes))
	if ptr == nil {
		return nil, fmt.Errorf("unable to allocate memory %d", numBytes)
	}
	return ptr, nil
}

func releaseMemory(ptr unsafe.Pointer) {
	C.release_memory(ptr)
}

func main() {
	// Parse command-line flags
	scheduleFile := flag.String("i", "", "path to the schedule file to execute, or empty for stdin")
	debugLevel := flag.Uint("d", 4, "debug level (0-6) where higher numbers have higher verbosity")
	cgroupsRootPath := flag.String("cgroups-root", "/sys/fs/cgroup", "path to the root of the cgroupsv1 hierarchy")
	flag.Parse()
	// Clamp debug level to valid range and set it
	if log.Level(*debugLevel) < log.PanicLevel {
		*debugLevel = uint(log.PanicLevel)
	}
	if log.Level(*debugLevel) > log.TraceLevel {
		*debugLevel = uint(log.TraceLevel)
	}
	log.SetLevel(log.Level(*debugLevel))
	// Load and execute the schedule
	schedule, err := loadSchedule(*scheduleFile, *cgroupsRootPath)
	if err != nil {
		log.Fatalf("failed to load schedule: %v", err)
	}
	if err := schedule.execute(); err != nil {
		log.Fatalf("failed to execute schedule: %v", err)
	}
}
