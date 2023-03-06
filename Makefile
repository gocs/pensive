
KEY := "sessionKey"

include .env
export

.PHONY: prod
prod:
	@GOOS="linux" GOARCH="amd64" CGO_ENABLED=0 go build -o app ./cmd/main.go
	@docker compose down && docker compose up --build -d


.PHONY: dev
dev:
	@GOOS="linux" GOARCH="amd64" CGO_ENABLED=0 go build -o ./test/app ./cmd/main.go
	@cd test && docker compose down && docker compose up --build -d



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