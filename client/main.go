package main

import (
	"log"

	pb "grpc/proto/helloworld"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "0.0.0.0:50051"
)

type Token struct {
	AppID     string
	AppSecret string
}

// RequireTransportSecurity implements credentials.PerRPCCredentials.
func (*Token) RequireTransportSecurity() bool {
	return false
}

// GetRequestMetadata 获取当前请求认证所需的元数据
func (t *Token) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"App-ID": t.AppID, "App-Secret": t.AppSecret}, nil
}

func main() {
	token := Token{
		AppID:     "1234567890",
		AppSecret: "GRPC-TOKEN",
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithPerRPCCredentials(&token))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := pb.NewGrpcClient(conn)
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Add: "TDsasd", Hour: 1, Value: 33000})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("code: %d, hash:%s", r.Code, r.Hash)
}
