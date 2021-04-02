TARGET:=hpc-workload-generator
SOURCES:=main.go helper.c helper.h schedule.go tasks.go cgroups/cgroups.go cgroups/cpuset.go

$(TARGET): $(SOURCES)
	go build -o $@

clean:
	rm -f $(TARGET)

.PHONY: clean
