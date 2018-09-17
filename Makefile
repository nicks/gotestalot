.PHONY: all build run

build:
	go build ./...

run:
	go run ./main.go -- github.com/windmilleng/tilt

test:
	go test ./...

install:
	go install github.com/nicks/gotestalot

watch:
	tilt up --watch main --debug
