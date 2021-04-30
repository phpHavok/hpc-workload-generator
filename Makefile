TARGET:=hpc-workload-generator
SOURCES:=main.go helper.c helper.h schedule.go tasks.go

$(TARGET): $(SOURCES)
	go build -o $@

clean:
	rm -f $(TARGET)

.PHONY: clean
