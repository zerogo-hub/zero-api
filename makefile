.PHONY: test
mod:
	go mod download
	go mod tidy
	go mod verify
	go mod vendor
test:
	go test ./router -coverprofile cover.out && go tool cover -html=cover.out -o cover.html