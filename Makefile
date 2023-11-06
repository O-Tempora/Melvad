BINARY_NAME = app
host = localhost
port = 5554

.PHONY:
	build \
	run \
	down \
	test \
	down-test \

help:
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

## build: builds project and creates binary file
build:
	go build -o $(BINARY_NAME) cmd/server/*.go

## run: runs the application. Specify host (redis host) and port (redis port) args
run: build 
	sudo docker run --name melvad-postgres -p 127.0.0.1:5553:5432/tcp -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=melvad -d postgres
	timeout 15s bash -c 'until sudo docker exec melvad-postgres pg_isready ; do sleep 1 ; done'
	sudo docker run --name melvad-redis -p $(port):6379 -d redis
	./$(BINARY_NAME) -host=$(host) -port=$(port)

## down: kills running containers with db and redis
down: 
	sudo docker stop melvad-redis && sudo docker stop melvad-postgres
	sudo docker rm melvad-redis && sudo docker rm melvad-postgres

## test: runs tests
test:
	sudo docker run --name melvad-postgres-test --pull=always -p 6553:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=melvad -d postgres
	timeout 15s bash -c 'until sudo docker exec melvad-postgres-test pg_isready ; do sleep 1 ; done'
	sudo docker run --name melvad-redis-test --pull=always -p 6554:6379 -d redis
	go test -v ./...
	$(MAKE) down-test

down-test: 
	sudo docker stop melvad-redis-test && sudo docker stop melvad-postgres-test
	sudo docker rm melvad-redis-test && sudo docker rm melvad-postgres-test
