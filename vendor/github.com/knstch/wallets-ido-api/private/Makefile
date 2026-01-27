PROTO_DIRS := $(sort $(dir $(wildcard **/*.proto)))
PROTO_FILES := $(wildcard *.proto)

.PHONY: all proto clean

all: proto

proto:
	@for file in $(PROTO_FILES); do \
		protoc --go_out=$(dir $$file) --go-grpc_out=$(dir $$file) -I. $$file; \
	done

clean:
	@find . -name "*.pb.go" -type f -delete