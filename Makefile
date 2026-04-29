.PHONY: fmt-check test vet smoke check

fmt-check:
	@test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)

test:
	go test ./...

vet:
	go vet ./...

smoke:
	sh scripts/smoke.sh

check: fmt-check test vet smoke
