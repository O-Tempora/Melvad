BINARY_NAME=server

.PHONY:
	build \
	run \
	down \

help:
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

## build: builds project and creates binary file
build:
	go build -o $(BINARY_NAME) cmd/api_server/*.go

## run: runs the application. Specify host (redis host) and port (redis port) args
run: build 
	sudo docker run --name melvad-redis -p 5554:6379 -d redis
	sudo docker run --name melvad-postgres -p 127.0.0.1:5553:5432/tcp -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=melvad -d postgres
	timeout 15s bash -c 'until sudo docker exec melvad-postgres pg_isready ; do sleep 2 ; done'
	./$(BINARY_NAME)
	@if [ -z $(port) ] && [ -z $(host) ]; then \
		./$(BINARY_NAME); \
	elif [ -z $(port) ]; then \
		/$(BINARY_NAME) -host=$(host); \
	elif [ -z $(host) ]; then \
		./$(BINARY_NAME) -port=$(port); \
	else \
		/$(BINARY_NAME) -host=$(host) -port=$(port); \
	fi

down: 
	sudo docker stop melvad-redis && sudo docker stop melvad-postgres
	sudo docker rm melvad-redis && sudo docker rm melvad-postgres
