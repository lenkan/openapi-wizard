build:
	mkdir -p out
	go build -o out/openapi ./cmd

test:
	go test ./...

demo:
	go run ./cmd --filename fixtures/openapi-example.yaml > fixtures/openapi-example-client.ts
