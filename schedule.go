package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/phpHavok/cgroups_exporter/cgroups"
	log "github.com/sirupsen/logrus"
)

type event struct {
	startTimeOffset time.Duration
	executable      executable
}

type schedule struct {
	processCgroups cgroups.Cgroups
	events         []event
}

func loadSchedule(filename string, cgroupsRootPath string) (schedule, error) {
	var schedule schedule
	// Load current cgroups before moving on
	processCgroups, err := cgroups.LoadProcessCgroups(os.Getpid(), cgroupsRootPath)
	if err != nil {
		return schedule, err
	}
	schedule.processCgroups = processCgroups
	// Load the schedule file, or default to stdin
	var scheduleFile *os.File
	if filename != "" {
		scheduleFile, err = os.Open(filename)
		if err != nil {
			return schedule, err
		}
		defer scheduleFile.Close()
	} else {
		scheduleFile = os.Stdin
		log.Info("Reading schedule of events from STDIN...")
	}
	csvReader := csv.NewReader(scheduleFile)
	csvReader.Comma = ','
	csvReader.Comment = '#'
	// Allow variable fields per record
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true
	csvLines, err := csvReader.ReadAll()
	if err != nil {
		return schedule, err
	}
	// Parse the CSV lines and sanity check along the way
	lastSecondsSeen := 0
	for _, csvLine := range csvLines {
		// At a minimum, we require a time offset and task name
		if len(csvLine) < 2 {
			return schedule, fmt.Errorf("invalid schedule line: %v", csvLine)
		}
		// The seconds time offset must be a valid integer and we require the
		// tasks to be listed in order
		seconds, err := strconv.Atoi(csvLine[0])
		if err != nil {
			return schedule, err
		}
		if seconds < lastSecondsSeen {
			return schedule, fmt.Errorf("task erroneously set to occur before a previous task: %v", csvLine)
		}
		lastSecondsSeen = seconds
		// Create a new executable from the given parameters
		executable, err := createExecutable(csvLine[1], csvLine[2:], processCgroups)
		if err != nil {
			return schedule, err
		}
		// Generate a new event on the schedule
		var event event
		event.startTimeOffset = time.Duration(seconds) * time.Second
		event.executable = executable
		schedule.events = append(schedule.events, event)
	}
	return schedule, nil
}

func createExecutable(taskName string, taskArgs []string, taskCgroups cgroups.Cgroups) (executable, error) {
	if taskName == "cpuload" {
		var task cpuLoadTask
		task.args = taskArgs
		if len(taskArgs) != 3 {
			return task, fmt.Errorf("task %s expected 3 args, got: %v", taskName, taskArgs)
		}
		// Load allocated CPUs from cgroup
		cpus, err := taskCgroups.Cpuset.GetCpus()
		if err != nil {
			return task, err
		}
		// Process cpuId
		cpuID, err := strconv.Atoi(taskArgs[0])
		if err != nil {
			return task, err
		}
		// Check if cpuID is within range
		if cpuID < 0 || cpuID >= len(cpus) {
			return task, fmt.Errorf("cpuID %d is out of range: (0-%d)", cpuID, len(cpus)-1)
		}
		// Map cpuID to physical CPU ID
		task.cpuID = cpus[cpuID]
		// Process pctLoad
		pctLoad, err := strconv.Atoi(taskArgs[1])
		if err != nil {
			return task, err
		}
		if pctLoad < 1 || pctLoad > 100 {
			return task, fmt.Errorf("pctLoad %d must be between 1 and 100 inclusive", pctLoad)
		}
		task.pctLoad = pctLoad
		// Process duration
		durationSecs, err := strconv.Atoi(taskArgs[2])
		if err != nil {
			return task, err
		}
		if durationSecs < 1 {
			return task, fmt.Errorf("duration %d (in seconds) must be positive", durationSecs)
		}
		task.duration = time.Duration(durationSecs) * time.Second
		return task, nil
	}
	return nil, fmt.Errorf("task %s is unknown", taskName)
}

// Execute the schedule to completion
func (s schedule) execute() error {
	cpus, err := s.processCgroups.Cpuset.GetCpus()
	if err != nil {
		return err
	}
	// We can execute up to one thread per CPU listed
	runtime.GOMAXPROCS(len(cpus))
	log.Infof("Input CPU indicies will be mapped to the cgroup CPUs: %v", cpus)
	numEvents := len(s.events)
	if numEvents < 1 {
		return fmt.Errorf("no events were scheduled")
	}
	// Actually run the schedule of events
	log.Info("Starting schedule of events...")
	statusCodes := make(chan int, numEvents)
	scheduleStart := time.Now()
	for _, event := range s.events {
		for time.Now().Sub(scheduleStart) < event.startTimeOffset {
			time.Sleep(500 * time.Millisecond)
		}
		log.Infof("Processing executable: ", event.executable)
		go func(executable executable, statusCodes chan int) {
			executable.execute(statusCodes)
		}(event.executable, statusCodes)
	}
	// Wait for processes to exit
	log.Info("Waiting for all processes to exit...")
	for i := 0; i < numEvents; i++ {
		<-statusCodes
	}
	log.Info("Schedule of events finished running to completion")
	return nil
}
