build:
	go build -o bin/github-activity
run: build
	./bin/github-activity
test:
	go test -v ./...
install: build
	go install -v
