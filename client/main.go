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

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := pb.NewGrpcClient(conn)
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Token: "123", Add: "TDsasd", Hour: 1, Value: 33000})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("code: %d, hash:%s", r.Code, r.Hash)
}
