TARGET:=hpc-workload-generator
CONTAINER:=hpc-workload-generator.sif

$(TARGET): main.go helper.c helper.h
	go build -o $@

$(CONTAINER): Singularity
	sudo singularity build $@ $<

container: $(CONTAINER)

clean:
	rm -f $(TARGET) $(CONTAINER)

.PHONY: clean container
