package cgroups

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	procCgroupIdxSubsystems = 1
	procCgroupIdxPath       = 2
)

// Cgroups represents a structure a cgroups across supported subsystems
type Cgroups struct {
	Cpuset cpuset
}

func readFile(root string, filename string) (string, error) {
	data, err := ioutil.ReadFile(filepath.Join(root, filename))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LoadProcessCgroups loads a structure containing all cgroups for a given process
func LoadProcessCgroups(pid int, cgroupsRootPath string) (Cgroups, error) {
	var cgroups Cgroups
	// Find and open the cgroup file for the process
	cgroupsPath := filepath.Join("/proc", strconv.Itoa(pid), "cgroup")
	cgroupsFile, err := os.Open(cgroupsPath)
	if err != nil {
		return cgroups, err
	}
	defer cgroupsFile.Close()
	// Load the cgroup file as CSV data
	csvReader := csv.NewReader(cgroupsFile)
	csvReader.Comma = ':'
	csvLines, err := csvReader.ReadAll()
	if err != nil {
		return cgroups, err
	}
	// Structure the CSV data into a map
	for _, csvLine := range csvLines {
		subsystems := strings.Split(csvLine[procCgroupIdxSubsystems], ",")
		for _, subsystem := range subsystems {
			cgroupAbsolutePath := filepath.Join(cgroupsRootPath, strings.TrimPrefix(subsystem, "name="), csvLine[procCgroupIdxPath])
			if _, err := os.Stat(cgroupAbsolutePath); os.IsNotExist(err) {
				return cgroups, fmt.Errorf("cgroup path doesn't exist: %s", cgroupAbsolutePath)
			}
			switch subsystem {
			case "cpuset":
				cgroups.Cpuset = cpuset(cgroupAbsolutePath)
			default:
				fmt.Println("Skipping unimplemented subsystem: ", subsystem)
			}
		}
	}
	return cgroups, nil
}
