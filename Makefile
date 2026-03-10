.PHONY: build run install clean

build:
	go build -o aitutor .

run:
	go run .

install:
	go install .

clean:
	rm -f aitutor

vet:
	go vet ./...
