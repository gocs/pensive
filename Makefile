
KEY := "sessionKey"

include .env
export

.PHONY: build up down

build:
	GOOS="linux" GOARCH="amd64" CGO_ENABLED=0 go build -o app ./cmd/main.go

up:
	docker compose up --build -d

down:
	docker compose down



# THIS DOESN'T WORK
# .PHONY: prod
# prod:
# 	@echo run this once
# 	@sudo apt install make docker-compose -y
# 	@wget https://golang.org/dl/go1.16.4.linux-amd64.tar.gz
# 	@sudo rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.4.linux-amd64.tar.gz
# 	@rm go1.16.4.linux-amd64.tar.gz
# 	@export PATH=$PATH:/usr/local/go/bin
# 	@source $HOME/.profile
# 	@echo ADD THE .env.prod file
# 	@echo ADD THE .env.prod file
# 	@echo ADD THE .env.prod file
# 	@echo ADD THE .env.prod file
# 	@echo ADD THE .env.prod file