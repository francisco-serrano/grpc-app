package main

import (
	"context"
	"github.com/francisco-serrano/grpc-app/decomposition/decompositionpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"strconv"
)

func doServerStreaming(client decompositionpb.PrimeNumberServiceClient, value int32) {
	log.Printf("attempting to perform Decomposition gRPC request...\n")

	req := &decompositionpb.DecompositionRequest{
		PrimeNumber: value,
	}

	responseStream, err := client.PrimeNumber(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumber RPC: %v", err)
	}

	for {
		msg, err := responseStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("response from PrimeNumber: %d\n", msg.GetPrimeElement())
	}
}

func main() {
	address := "localhost:50051"

	if len(os.Args) != 2 {
		log.Fatalf("only one integer argument required")
	}

	value, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("could not parse cli argument, is it integer?")
	}

	log.Printf("starting gRPC client at %s with param %d...\n", address, value)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to %s: %v", address, err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("could not close connection properly: %v", err)
		}
	}()

	client := decompositionpb.NewPrimeNumberServiceClient(conn)

	doServerStreaming(client, int32(value))
}
