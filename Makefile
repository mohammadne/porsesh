unit-tests:
	go test -v ./internal/... -covermode=atomic -coverprofile=coverage.out

load-test:
	k6 run tests/stress/k6.js
