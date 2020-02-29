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

func doBidirectionalStreaming(c greetpb.GreetServiceClient) {
	log.Printf("starting to perform a bidirectional streaming RPC...")

	// we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream: %v", err)
	}

	waitc := make(chan struct{})

	requests := []*greetpb.GreetEveryoneRequest{
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

	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			log.Printf("sending message %+v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}

		stream.CloseSend()
	}()

	// we receive a bunch of message from the client (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving: %v", err)
			}

			log.Printf("received: %v\n", res.GetResult())
		}

		close(waitc)
	}()

	// block until everything is done
	<-waitc
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

	//doClientStreaming(c)

	doBidirectionalStreaming(c)
}
