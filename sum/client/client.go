package main

import (
	"context"
	"github.com/francisco-serrano/grpc-app/sum/sumpb"
	"google.golang.org/grpc"
	"log"
)

func doUnary(c sumpb.SumServiceClient) {
	request := &sumpb.SumRequest{
		FirstElement:  10,
		SecondElement: 5,
	}

	log.Printf("Starting to do gRPC request with values %v and %v", request.FirstElement, request.SecondElement)

	response, err := c.Sum(context.Background(), request)
	if err != nil {
		log.Fatalf("could not perform gRPC sum: %v", err)
	}

	log.Printf("sum over gRPC successful, result: %v", response.Result)
}

func main() {
	log.Printf("starting gRPC client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()

	c := sumpb.NewSumServiceClient(conn)

	doUnary(c)
}
