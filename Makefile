EXE_NAME := wafer

.PHONY: all generate build run clean

all: build

builtins:
	python3 generate_builtins.py

build: builtins
	cd src && go build -o ../$(EXE_NAME)

run: build
	./$(EXE_NAME) $(ARGS)

clean:
	rm -f $(EXE_NAME)
	rm -f src/builtins_generated.go