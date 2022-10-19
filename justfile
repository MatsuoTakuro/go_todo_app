run:
  go run . 18080

test:
  go test -timeout 30s -v

build-dev:
  docker compose build --no-cache

run-dev:
  docker compose up