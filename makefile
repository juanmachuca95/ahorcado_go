gen-grpc-go:
	protoc -I=protos --go_out=. --go-grpc_out=. protos/proto.proto

clean-grpc-go:
	rm -rf generated/*.pb.go

gen-grpc-gateway:
	protoc -I ./protos \
	--go_out ./protos --go_opt paths=source_relative \
	--go-grpc_out ./protos --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./protos --grpc-gateway_opt paths=source_relative \
	./protos/ahorcado/proto.proto \
	./protos/auth/proto.proto
	
clean-grpc-gateway:
	rm -rf protos/ahorcado/*.pb.gw.go

compile-go-js:
	gopherjs build --minify cmd/clients/js/client.go -o cmd/clients/js/html/index.js

setup-site-github:
	cp cmd/clients/js/html/index.html ./docs 
	cp cmd/clients/js/html/index.js ./docs 

clean-site-github:
	rm -rf docs/*

up-server:
	go run server/server.go