go get github.com/cespare/reflex && \
reflex -r '\.go|json$' -s -- sh -c 'go build && ./noah  --config=./example/postgres.yml --dir=./example/postgres -D';