package main

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc"

	pb "github.com/dulei0726/grpc-demo/server_stream/proto"
)

func SayHello(client pb.GreeterClient) error {
	log.Println("client before")
	stream, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "dulei"})
	if err != nil {
		return err
	}
	log.Println("client after")

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("client resp: %s", resp.Message)
	}
	return nil
}

// HelloInterceptor returns a client interceptor example
func HelloInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		log.Println("Hello before")
		stream, err := streamer(ctx, desc, cc, method, opts...)
		log.Println("Hello after")
		return stream, err
	}
}

// HiInterceptor returns a client interceptor example
func HiInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		log.Println("Hi before")
		stream, err := streamer(ctx, desc, cc, method, opts...)
		log.Println("Hi after")
		return stream, err
	}
}

func main() {
	conn, err := grpc.Dial(
		":8080",
		grpc.WithInsecure(),
		grpc.WithChainStreamInterceptor(HelloInterceptor(), HiInterceptor()),
	)
	if err != nil {
		log.Fatal("Dial err:", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	err = SayHello(client)
	if err != nil {
		log.Println("SayHello err", err)
	}
}
