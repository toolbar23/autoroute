all: build

build: .
	go build -o ./target/autoroute ./cmd/autoroute

run: build
	./target/autoroute