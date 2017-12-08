.PHONY: compile
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go

# If $GOPATH/bin/protoc-gen-go does not exist, we'll run this command to install it.
$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

contacts/contacts.pb.go: contacts/contacts.proto | $(PROTOC_GEN_GO)
	protoc -I contacts/ contacts/contacts.proto --go_out=plugins=grpc:contacts

# This is a "phony" target - an alias for the above command, so "make compile"
# still works.
compile: contacts/contacts.pb.go
