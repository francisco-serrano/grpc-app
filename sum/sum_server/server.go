package main

import (
	"context"
	"github.com/francisco-serrano/grpc-app/sum/sumpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *sumpb.SumRequest) (*sumpb.SumResponse, error) {
	log.Printf("Sum function was invoked with %v", req)

	result := req.GetFirstElement() + req.GetSecondElement()

	log.Printf("Sum function result was %d", result)

	return &sumpb.SumResponse{
		Result: result,
	}, nil
}

func main() {
	log.Printf("Attempting to start sum grpc s")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	sumpb.RegisterSumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
