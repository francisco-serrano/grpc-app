package main

import (
	"context"
	"github.com/francisco-serrano/grpc-app/average/averagepb"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
)

func doClientStreaming(client averagepb.AverageServiceClient, values []int32) {
	log.Printf("attempting to perform Average RPC...")

	stream, err := client.Average(context.Background())
	if err != nil {
		log.Fatalf("could not obtain stream: %+v", err)
	}

	for _, v := range values {
		log.Printf("sending %d value to server...", v)

		if err := stream.Send(&averagepb.AverageRequest{
			Number: v,
		}); err != nil {
			log.Fatalf("could not send message with number %d: %+v", v, err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while trying to close stream: %+v", err)
	}

	log.Printf("obtained response from Average RPC: %+v\n", res)
}

func main() {
	address := "localhost:50051"

	if len(os.Args) == 1 {
		log.Fatalf("elements required as cli arguments in order to obtain average")
	}

	var values []int32
	for _, arg := range os.Args[1:] {
		v, err := strconv.Atoi(arg)
		if err != nil {
			log.Fatalf("average only supports integer values")
		}

		values = append(values, int32(v))
	}

	log.Printf("attempting to start gRPC client at %s\n", address)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to gRPC server: %+v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("could not close server connection: %+v", err)
		}
	}()

	client := averagepb.NewAverageServiceClient(conn)

	doClientStreaming(client, values)
}
