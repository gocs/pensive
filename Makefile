
KEY := "sessionKey"

.PHONY: reload
reload:
	@docker-compose down
	@docker-compose up --build -d


.PHONY: run
run:
	@go run cmd/main.go -session-key=$(KEY)