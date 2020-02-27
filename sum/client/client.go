package main

import (
	"context"
	"github.com/francisco-serrano/grpc-app/sum/sumpb"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
)

func doUnary(c sumpb.SumServiceClient, a, b int32) {
	request := &sumpb.SumRequest{
		FirstElement:  a,
		SecondElement: b,
	}

	log.Printf("Starting to do gRPC request with values %v and %v", request.FirstElement, request.SecondElement)

	response, err := c.Sum(context.Background(), request)
	if err != nil {
		log.Fatalf("could not perform gRPC sum: %v", err)
	}

	log.Printf("sum over gRPC successful, result: %v", response.Result)
}

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		log.Fatalf("need to cli operands to perform sum")
	}

	a, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("only integer operands supported")
	}

	b, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("only integer operands supported")
	}

	log.Printf("starting gRPC client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()

	c := sumpb.NewSumServiceClient(conn)

	doUnary(c, int32(a), int32(b))
}
