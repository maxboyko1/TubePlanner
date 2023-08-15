.PHONY = all clean

EXECS = tubeplanner

all: $(EXECS)

tubeplanner: tubeplanner.go transitdata.go
	go build -o $@ $(wildcard *.go)

clean:
	@rm -f $(EXECS)
