build_bot:
	go build -o bot cmd/main.go

build_api:
	go build -o api api/*.go