package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	pb "github.com/dulei0726/grpc-demo/client_stream/proto"
)

func SayHello(client pb.GreeterClient) error {
	stream, err := client.SayHello(context.Background())
	if err != nil {
		return err
	}
	for i := 0; i <= 6; i++ {
		err = stream.Send(&pb.HelloRequest{Name: fmt.Sprintf("dulei %d", i)})
		if err != nil {
			return err
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("client resp: %s", resp.Message)
	return nil
}

func HelloInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		return clientStream, err
	}
}

func main() {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure(), grpc.WithStreamInterceptor(HelloInterceptor()))
	if err != nil {
		log.Fatal("grpc.Dial err:", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	err = SayHello(client)
	if err != nil {
		log.Println("SayHello err:", err)
	}
}
