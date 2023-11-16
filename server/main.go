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
	log.Printf("ADD=%s, \nHour=%d, \nValue=%d", in.Add, in.Hour, in.Value)
	return &pb.HelloReply{Code: 200, Hash: "hash" + in.Add}, nil
}

func CustomHeaderMatcher(key string) (string, bool) {
	switch key {
	case "App-Id":
		return key, true
	case "App-Secret":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func middleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Println("进入拦截器验证")
	token, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("无Token认证信息")
		return nil, status.Errorf(codes.Unauthenticated, "无Token认证信息")
	}
	var (
		AppID     string
		AppSecret string
	)
	if val, ok := token["app-id"]; ok {
		AppID = val[0]
	}
	if val, ok := token["app-secret"]; ok {
		AppSecret = val[0]
	}
	log.Printf("metadata: %v \nAppSecret=%s, \nAppID=%s", token, AppSecret, AppID)
	if len(AppSecret) == 0 || len(AppID) == 0 {
		log.Println("AppID或AppSecret验证失败")
		return nil, status.Errorf(codes.Unauthenticated, "AppID 或 AppSecret 验证失败")
	}
	log.Println("验证成功, 进入下一步")
	return handler(ctx, req)
}

func main() {
	// Create a gRPC server object
	s := grpc.NewServer(grpc.UnaryInterceptor(middleware))
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
	gwmux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(CustomHeaderMatcher))
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
