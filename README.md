# gRPC
### 同时提供 gRPC 和 http 服务的示例

	protoc --proto_path=./proto \
	    --go_out=./proto --go_opt=paths=source_relative \
	    --go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
	   ./proto/helloworld/hello_world.proto
	
	
	
	protoc --proto_path=./proto \
	    --go_out=./proto --go_opt=paths=source_relative \
	    --go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
	    --grpc-gateway_out=proto --grpc-gateway_opt=paths=source_relative \
	   ./proto/helloworld/hello_world.proto

#### CURL命令用来测试

	curl --location 'http://localhost:8080/v1/grpc/sayhello' \
	--header 'App-secret: 123456' \
	--header 'App-id: kkkkkkang1212' \
	--data '{"value":"210920192",
	"hour":1,
	"add":"TBDASJDAKSDKLSKAD",
	"token":"1232121"}'
