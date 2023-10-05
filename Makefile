compile:
	protoc api/*.proto \
	--go_out=. \
	--go_opt=paths=import \
	--go-grpc_out=. --go-grpc_opt=paths=import \
	--proto_path=.