BINARY_NAME=meepShopTest

build:
	go build -o ${BINARY_NAME} main.go

run: build
	./${BINARY_NAME}

docker-build-server:
	docker build -t my-golang-app .

docker-up-server:
	docker run -it -p 8080:8080 --name testName my-golang-app

docker-up-db:
	docker pull postgres:14.9
	docker run -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=test -e POSTGRES_USER=test -e POSTGRES_DB=meepShopTest postgres:14.9
	docker cp ./meepShopTest.sql postgres:/file.sql
	sleep 2;
	docker exec -u postgres postgres psql -U test meepShopTest -f /file.sql

docker-up:
	docker compose up -d

docker-down:
	docker compose down

# mock: install-mocks
# 	mockgen --source=./internal/service/user_test.go --destination ./internal/service/user_mock.go --package repository

# install-mocks: # could be replaced by local bin to avoid different version
# 	@go get github.com/golang/mock/gomock
# 	@go install github.com/golang/mock/mockgen@v1.6.0

test:
	go test ./...