install:
	go install -v

build:
	go build -v ./...

lint:
	gometalinter --config .linter.conf

test:
	go test --cover -v $(go list ./... | grep -v /vendor/)

deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

dev-deps: deps
	go get -u github.com/smartystreets/goconvey/convey
	go get github.com/alecthomas/gometalinter
	gometalinter --install

clean:
	go clean

dist-clean:
	rm -rf pkg src bin
