install:
	go install -v

build:
	go build -v ./...

lint:
	gometalinter --config .linter.conf

test:
	go test -v ./... --cover

deps:
	go get -u gopkg.in/r3labs/graph.v2
	go get -u github.com/tidwall/gjson

dev-deps: deps
	go get -u github.com/smartystreets/goconvey/convey
	go get github.com/alecthomas/gometalinter
	gometalinter --install

clean:
	go clean

dist-clean:
	rm -rf pkg src bin
