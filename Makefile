SRC_PATH?=./cmd/theredshirts-message
APP_NAME?=theredshirts-message
DOCKER_PATH?=./build/Dockerfile

app.build:
	go mod download
	go build -o $(APP_NAME) $(SRC_PATH)

docker.build:
	docker build -t beancodede/$(APP_NAME):latest -f $(DOCKER_PATH) .