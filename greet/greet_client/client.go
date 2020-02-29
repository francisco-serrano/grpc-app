package main

import (
	"context"
	"fmt"
	"github.com/francisco-serrano/grpc-app/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func doUnary(c greetpb.GreetServiceClient) {
	log.Printf("Starting to do a Unary RPC...")

	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Francisco",
			LastName:  "Serrano",
		},
	}

	response, err := c.Greet(context.Background(), request)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", response.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Francisco",
			LastName:  "Serrano",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("response from GreetManyTimes: %s", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	log.Printf("Starting to do a Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Francisco",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Pepe",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Juan",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Jose",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Raul",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	for _, req := range requests {
		log.Printf("Sending req: %+v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}

	log.Printf("LongGreet response: %+v\n", res)
}

func main() {
	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	//doUnary(c)

	//doServerStreaming(c)

	doClientStreaming(c)
}
