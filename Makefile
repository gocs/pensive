
KEY := "sessionKey"

.PHONY: reload
reload:
	@GOOS="linux" GOARCH="amd64" CGO_ENABLED=0 go build -o app ./cmd/main.go
	@docker-compose down && docker-compose up --build -d


.PHONY: run
run:
	@go run cmd/main.go -session-key=$(KEY)
