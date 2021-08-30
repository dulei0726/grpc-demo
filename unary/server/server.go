package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/dulei0726/grpc-demo/unary/proto"
)

// GreeterServer implement pb.GreeterServer interface
type GreeterServer struct{
	pb.UnimplementedGreeterServer
}

func (GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println(ctx.Deadline())
	log.Println("md: ", md)
	return &pb.HelloReply{
		Message: "Hello, " + r.Name,
	}, nil
}

// HelloInterceptor returns a server interceptor example
func HelloInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		r := req.(*pb.HelloRequest)
		log.Println("Hello before", r.Name)
		resp, err = handler(ctx, req)
		log.Println("Hello after", r.Name)
		return
	}
}

// HiInterceptor returns a server interceptor example
func HiInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		r := req.(*pb.HelloRequest)
		log.Println("Hi before", r.Name)
		resp, err = handler(ctx, req)
		log.Println("Hi after", r.Name)
		return
	}
}

func main() {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(HelloInterceptor(), HiInterceptor()),
	)
	pb.RegisterGreeterServer(s, GreeterServer{})
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("net.Listen err:", err)
	}
	err = s.Serve(lis)
	if err != nil {
		log.Fatal("Serve err:", err)
	}
}
