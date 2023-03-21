compose:
	docker compose down
	docker compose up -d --build

install-dependencies:
	# format tools
	go install github.com/daixiang0/gci@latest
	go install github.com/momaek/formattag@latest
	go install mvdan.cc/gofumpt@latest

.PHONY: format
format:
	# Format
	(find . -name "*.go" -exec go fmt {} \;)
	(find . -name "*.go" -exec gofumpt -l -w {} \;)
	(find . -name "*.go" -exec formattag -file {} \;)
	(find . -name "*.go" -exec gci write {} \;)

start:
	docker compose up -d