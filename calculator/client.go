package main

import (
	"context"
	"fmt"
	"github.com/francisco-serrano/grpc-app/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"sync"
)

func doUnarySum(c calculatorpb.CalculatorServiceClient, a, b int32) {
	request := &calculatorpb.SumRequest{
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

func doServerStreamingDecomposition(client calculatorpb.CalculatorServiceClient, value int32) {
	log.Printf("attempting to perform Decomposition gRPC request...\n")

	req := &calculatorpb.DecompositionRequest{
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

func doClientStreamingAverage(client calculatorpb.CalculatorServiceClient, values []int32) {
	log.Printf("attempting to perform Average RPC...")

	stream, err := client.Average(context.Background())
	if err != nil {
		log.Fatalf("could not obtain stream: %+v", err)
	}

	for _, v := range values {
		log.Printf("sending %d value to server...", v)

		if err := stream.Send(&calculatorpb.AverageRequest{
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

func doBidirectionalStreamingMaximum(c calculatorpb.CalculatorServiceClient, values []int32) {
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
			req := &calculatorpb.MaximumRequest{
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

func doErrorUnary(c calculatorpb.CalculatorServiceClient, number int32) {
	log.Printf("starting to do a SquareRoot Unary RPC...")

	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: number,
	})
	if err != nil {
		resp, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			log.Printf("error message from server: %v\n", resp.Message())
			fmt.Println(resp.Code())
			if resp.Code() == codes.InvalidArgument {
				log.Printf("we probably sent a negative number\n")
			}
		} else {
			log.Fatalf("big error calling SquareRoor: %v", err)
		}
	}

	log.Printf("result of square root of %v: %v", number, res.GetNumberRoot())
}

func main() {
	address := "localhost:50051"

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error while trying to connect to server on %d: %v", address, err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("error while trying to close connection with server: %v", err)
		}
	}()

	c := calculatorpb.NewCalculatorServiceClient(conn)

	//doUnarySum(c, 20, 31)
	//
	//doServerStreamingDecomposition(c, 343)
	//
	//doClientStreamingAverage(c, []int32{2, 5, 8, 3, 1, 6, 72, 14, 5})
	//
	//doBidirectionalStreamingMaximum(c, []int32{2, 5, 3, 20, 17, 40, 35, 50})

	doErrorUnary(c, 10)
	doErrorUnary(c, -2)
}
