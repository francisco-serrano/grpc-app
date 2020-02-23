package main

import (
	"context"
	"fmt"
	"github.com/francisco-serrano/grpc-app/greet/greetpb"
	"google.golang.org/grpc"
	"log"
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

func main() {
	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	doUnary(c)
}
