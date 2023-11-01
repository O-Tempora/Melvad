BINARY_NAME=server

.PHONY:
	build \
	run \

help:
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

## build: builds project and creates binary file
build:
	go build -o $(BINARY_NAME) cmd/api_server/*.go

## run: runs the application. Specify host (redis host) and port (redis port) args
run: build 
	@if [ -z $(port) ] && [ -z $(host) ]; then \
		./$(BINARY_NAME); \
	elif [ -z $(port) ]; then \
		/$(BINARY_NAME) -host=$(host); \
	elif [ -z $(host) ]; then \
		./$(BINARY_NAME) -port=$(port); \
	else \
		/$(BINARY_NAME) -host=$(host) -port=$(port); \
	fi
