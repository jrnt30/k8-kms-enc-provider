default: ; @true

# TODO
# For each .proto file, ensure there is a correpsonding .db.go file present
.PHONY: protos
protos:
	protoc --go_out plugins=grpc:generated/ --proto_path proto/ proto/*.proto

docker-build: 
	docker build -t kms-server . 

