package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	pb "github.com/dulei0726/grpc-demo/both_stream/proto"
)

func SayHello(client pb.GreeterClient) error {
	stream, err := client.SayHello(context.Background())
	if err != nil {
		return err
	}

	var sendFunc = func() error {
		for i := 0; i <= 6; i++ {
			err := stream.Send(&pb.HelloRequest{Name: fmt.Sprintf("dulei %d", i)})
			if err != nil {
				return err
			}
		}
		return stream.CloseSend()
	}

	var recvFunc = func() error {
		for {
			r, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			log.Printf("client recv: %s\n", r.Message)
		}
	}

	var eg errgroup.Group
	eg.Go(sendFunc)
	eg.Go(recvFunc)

	return eg.Wait()
}

func main() {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("grpc.Dial err:", err)
		return
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	err = SayHello(client)
	if err != nil {
		log.Fatal("SayHello err:", err)
	}
}
