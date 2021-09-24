generate:
	buf mod update -v
	buf build -v
	buf generate -v
	go mod tidy

build:
	go build -o ./bin/server -v ./cmd


up:
	docker-compose build
	docker-compose up -d pg server 

down:
	docker-compose down pg server

stop:
	docker-compose stop pg server

start:
	docker-compose start pg server