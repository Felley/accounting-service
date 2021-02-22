.PHONY: protos

protos:
	protoc -I protos/ --go-grpc_out=protos/accounting --go_out=protos/accounting protos/accounting.proto