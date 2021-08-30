package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/dulei0726/grpc-demo/server_stream/proto"
)

type GreeterServer struct{
	pb.UnimplementedGreeterServer
}

func (GreeterServer) SayHello(r *pb.HelloRequest, stream pb.Greeter_SayHelloServer) error {
	for i := 0; i <= 6; i++ {
		err := stream.Send(&pb.HelloReply{Message: fmt.Sprintf("Hello, %s, %d", r.Name, i)})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	s := grpc.NewServer()
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
