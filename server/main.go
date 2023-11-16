package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	pb "grpc/proto/helloworld"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var port = flag.Int("port", 50051, "the port to serve on")
var restful = flag.Int("restful", 8080, "the port to restful serve on")

type server struct {
	pb.UnimplementedGrpcServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	token, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "无Token认证信息")
	}
	var (
		App_ID     string
		App_Secret string
	)
	if val, ok := token["app-id"]; ok {
		App_ID = val[0]
	}
	if val, ok := token["app-secret"]; ok {
		App_Secret = val[0]
	}

	log.Printf("metadata: %v, App_Secret=%s, App_ID=%s", token, App_Secret, App_ID)
	return &pb.HelloReply{Code: 200, Hash: "hash" + in.Add}, nil
}

func CustomMatcher(key string) (string, bool) {
	switch key {
	case "App-Id":
		return key, true
	case "App-Secret":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func main() {
	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Grpc service to the server
	pb.RegisterGrpcServer(s, &server{})
	// Serve gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Serving gRPC on 0.0.0.0" + fmt.Sprintf(":%d", *port))
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// 2. 启动 HTTP 服务
	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	// 3 自定义 HTTP headers 规则
	gwmux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(CustomMatcher))
	// Register Grpc
	err = pb.RegisterGrpcHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *restful),
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0" + fmt.Sprintf(":%d", *restful))
	log.Fatalln(gwServer.ListenAndServe())
}
