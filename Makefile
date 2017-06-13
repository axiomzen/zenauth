.PHONY: ci install test.unit test.integrate build_docs build_protobuf
ci: install test.unit test.integrate
install:
	zest build
	zest bundle
test.unit:
	zest test
test.integrate:
	zest integrate
build_docs:
	swagger-codegen generate -l html -i swagger.yml
	mv index.html docs.html
build_protobuf:
	protoc -I protobuf/ protobuf/*.proto --go_out=plugins=grpc:protobuf
