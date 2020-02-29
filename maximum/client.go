package main

import (
	"context"
	"github.com/francisco-serrano/grpc-app/maximum/maximumpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
)

func doBidirectionalStreaming(c maximumpb.MaximumServiceClient, values []int32) {
	log.Printf("attempting to perform Maximum RPC call...")

	stream, err := c.Maximum(context.Background())
	if err != nil {
		log.Fatalf("error while trying to obtain stream: %v", err)
	}

	wg := sync.WaitGroup{}

	// send values over the streaming
	wg.Add(1)
	go func() {
		for _, v := range values {
			req := &maximumpb.MaximumRequest{
				Value: v,
			}

			log.Printf("sending %+v to the server\n", req)

			if err := stream.Send(req); err != nil {
				log.Fatalf("error while sending message to server: %v", err)
			}
		}

		if err := stream.CloseSend(); err != nil {
			log.Fatalf("error while attempting to send stream to server: %v", err)
		}

		wg.Done()
	}()

	// receive values from the streaming
	wg.Add(1)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving message to server: %v", err)
			}

			log.Printf("received %+v from server\n", res)
		}

		wg.Done()
	}()

	wg.Wait()

	log.Printf("bidirectional streaming performed successfully\n")
}

func main() {
	address := "localhost:50051"

	if len(os.Args) == 1 {
		log.Fatalf("need some cli integer values to perform streaming")
	}

	var values []int32

	for _, arg := range os.Args[1:] {
		v, err := strconv.Atoi(arg)
		if err != nil {
			log.Fatalf("error while trying to parse cli arguments: %v", err)
		}

		values = append(values, int32(v))
	}

	log.Printf("attempting to connect to gRPC server at %s\n", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error while trying to connect to gRPC server: %v", err)
	}

	c := maximumpb.NewMaximumServiceClient(conn)

	doBidirectionalStreaming(c, values)
}
