package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/dulei0726/grpc-demo/unary/proto"
)

func SayHello(client pb.GreeterClient) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	ctx = metadata.AppendToOutgoingContext(ctx, "key", "value")
	log.Println("client before")
	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "dulei"})
	if err != nil {
		return err
	}
	log.Printf("client after: %s", resp.Message)
	return nil
}

// HelloInterceptor returns a client interceptor example
func HelloInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		resp := reply.(*pb.HelloReply)
		log.Println("Hello before", resp.Message)
		err := invoker(ctx, method, req, reply, cc, opts...)
		log.Println("Hello after", resp.Message)
		return err
	}
}

// HiInterceptor returns a client interceptor example
func HiInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		resp := reply.(*pb.HelloReply)
		log.Println("Hi before", resp.Message)
		err := invoker(ctx, method, req, reply, cc, opts...)
		log.Println("Hi after", resp.Message)
		return err
	}
}

func main() {
	conn, err := grpc.Dial(
		":8080",
		grpc.WithInsecure(),
		grpc.WithChainUnaryInterceptor(HelloInterceptor(), HiInterceptor()),
	)
	if err != nil {
		log.Fatal("Dial err:", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	err = SayHello(client)
	if err != nil {
		log.Println("SayHello err:", err)
	}
}
