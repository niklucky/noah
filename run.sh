	go get github.com/cespare/reflex && \
	reflex -r '\.go|json$' -s -- sh -c 'go build && ./noah  --config=./example/test.json --dir=./example/migrations -D';