.PHONY: hatch install test start run
hatch: install test
install:
	go get github.com/tools/godep
	godep restore
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' .
#	was axiomzen/context-engine:accessor
	
test:
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega
	go install
	ginkgo --progress -race test/integration
start:

run: | install start

build_docs:
	swagger-codegen generate -l html -i swagger.yml
	mv index.html docs.html

build_protobuf:
	protoc -I protobuf/ protobuf/*.proto --go_out=plugins=grpc:protobuf
