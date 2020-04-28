build_db:
	docker build -t shorts_db:1.0 .
	docker run --publish 5432:5432 --detach --name shorts shorts_db:1.0

start_db:
	docker start shorts

docs:
	set GO111MODULE=off
	swagger generate spec -o ./swagger.yaml --scan-models

test:
	go test

run:
	go run main.go