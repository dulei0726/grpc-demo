package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/dulei0726/grpc-demo/both_stream/proto"
)

type GreeterServer struct{
	pb.UnimplementedGreeterServer
}

func (GreeterServer) SayHello(stream pb.Greeter_SayHelloServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("server recv: %s\n", r.Name)

		err = stream.Send(&pb.HelloReply{Message: fmt.Sprintf("hello %s", r.Name)})
		if err != nil {
			return err
		}
		err = stream.Send(&pb.HelloReply{Message: fmt.Sprintf("hi %s", r.Name)})
		if err != nil {
			return err
		}
	}
}

func main() {
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, GreeterServer{})
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("net.Listen err:", err)
		return
	}
	err = s.Serve(lis)
	if err != nil {
		log.Fatal("Serve err:", err)
	}
}
