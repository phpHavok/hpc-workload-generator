TARGET:=hpc-workload-generator
CONTAINER:=hpc-workload-generator.sif
SOURCES:=main.go helper.c helper.h schedule.go tasks.go cgroups/cgroups.go cgroups/cpuset.go

$(TARGET): $(SOURCES)
	go build -o $@

$(CONTAINER): Singularity
	sudo singularity build $@ $<

container: $(CONTAINER)

clean:
	rm -f $(TARGET) $(CONTAINER)

.PHONY: clean container
