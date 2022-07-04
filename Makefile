GRPC_EXTERNAL=./external/grpc/wbroker
PROTO_FILES=$(shell find . -name "*.proto" | sort)
GRPC_PORT=24005

.PHONY: run
run: docker_build docker_run

docker_build:
	docker build -t wbroker .
docker_run:
	docker run -d -p $(GRPC_PORT):$(GRPC_PORT) wbroker

compile_proto: compile_proto_code compile_proto_doc

compile_proto_code:
	@if [ ! -z "$(PROTO_FILES)" ] && [ -d $(GRPC_EXTERNAL) ] ; then \
		echo "Compiling proto files"; \
		protoc \
		-I=$(GRPC_EXTERNAL) \
		--go_out=$(GRPC_EXTERNAL) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GRPC_EXTERNAL) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES); \
	else \
	  echo "No proto files found"; \
	fi

compile_proto_doc:
	@if [ ! -z "$(PROTO_FILES)" ] && [ -d $(GRPC_EXTERNAL) ] ; then \
		echo "Generating gRPC documentation"; \
		protoc -I=$(GRPC_EXTERNAL) --doc_out=./docs --doc_opt=markdown,grpc.md $(PROTO_FILES); \
        protoc -I=$(GRPC_EXTERNAL) --doc_out=./docs --doc_opt=html,grpc.html $(PROTO_FILES); \
    fi