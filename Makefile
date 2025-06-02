lint:
	golangci-lint run

build:
	go build -o tui-cli

run: 
	./tui-cli