BINARY_NAME=server
REDIS_HOST=localhost
REDIS_PORT=5554

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
ifneq ($(strip $(host)),)
	REDIS_HOST = $(host)
endif
ifneq ($(strip $(port)),)
	REDIS_PORT = $(port)
endif
	sudo docker run --name melvad-postgres -p 127.0.0.1:5553:5432/tcp -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=melvad -d postgres
	timeout 15s bash -c 'until sudo docker exec melvad-postgres pg_isready ; do sleep 1 ; done'
	sudo docker run --name melvad-redis -p $(REDIS_PORT):6379 -d redis
	./$(BINARY_NAME) -host=$(REDIS_HOST) -port=$(REDIS_PORT)

down: 
	sudo docker stop melvad-redis && sudo docker stop melvad-postgres
	sudo docker rm melvad-redis && sudo docker rm melvad-postgres
