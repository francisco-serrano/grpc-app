package main

import (
	"context"
	"fmt"
	"github.com/francisco-serrano/grpc-app/calculator/calculatorpb"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"
)

type server struct {
}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Printf("Sum function was invoked with %v", req)

	result := req.GetFirstElement() + req.GetSecondElement()

	log.Printf("Sum function result was %d", result)

	return &calculatorpb.SumResponse{
		Result: result,
	}, nil
}

func (s *server) PrimeNumber(req *calculatorpb.DecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberServer) error {
	log.Printf("received grpc request: %+v\n", req)

	N := req.GetPrimeNumber()
	k := int32(2)

	for N > 1 {
		if N%k == 0 {
			res := &calculatorpb.DecompositionResponse{
				PrimeElement: k,
			}

			log.Printf("attempting to send %+v\n", res)

			if err := stream.Send(res); err != nil {
				return err
			}

			log.Printf("sent %+v successfully\n", res)

			N /= k
		} else {
			k += 1
		}
	}

	log.Printf("stream successfully finalized\n")

	return nil
}

func (s *server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	log.Printf("Average function was invoked with a streaming request\n")

	var sum int32 = 0
	var amountValues int32 = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Result: &wrappers.DoubleValue{
					Value: float64(sum) / float64(amountValues),
				},
			})
		}

		if err != nil {
			log.Fatalf("error while reading client stream: %+v", err)
		}

		sum += req.GetNumber()
		amountValues += 1
	}
}

func (s *server) Maximum(stream calculatorpb.CalculatorService_MaximumServer) error {
	log.Printf("Maximum function was invoked with a streaming request\n")

	currentMax := int32(-1)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}

		log.Printf("received message from client stream: %+v\n", req)

		value := req.GetValue()

		if value > currentMax {
			currentMax = value

			res := &calculatorpb.MaximumResponse{
				MaxValue: currentMax,
			}

			log.Printf("sending message to client stream %+v\n", res)

			if err := stream.Send(res); err != nil {
				log.Fatalf("error while sending message to client stream: %+v", err)
			}
		}
	}
}

func (s *server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	log.Printf("SquareRoot invoked over RPC with params: %+v\n", req)

	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("received a negative number: %v", number),
		)
	}

	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	address := "localhost:50051"

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("could not listen on %s: %d", address, err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	reflection.Register(s)

	log.Printf("running server on %s...", address)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("error while trying to serve application: %v", err)
	}
}
