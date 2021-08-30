package main

import (
	"fmt"
	"io"
	"net"

	"google.golang.org/grpc"

	pb "github.com/dulei0726/grpc-demo/client_stream/proto"
)

type GreeterServer struct{
	pb.UnimplementedGreeterServer
}

func (GreeterServer) SayHello(stream pb.Greeter_SayHelloServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			message := &pb.HelloReply{Message: "server recv end"}
			return stream.SendAndClose(message)
		}
		if err != nil {
			return err
		}
		fmt.Printf("server recv: %s\n", r.Name)
	}
}

func main() {
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, GreeterServer{})
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	err = s.Serve(lis)
	if err != nil {
		fmt.Println("Serve err:", err)
	}
}
