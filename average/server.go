package main

import (
	"github.com/francisco-serrano/grpc-app/average/averagepb"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct {
}

func (s *server) Average(stream averagepb.AverageService_AverageServer) error {
	log.Printf("Average function was invoked with a streaming request\n")

	var sum int32 = 0
	var amountValues int32 = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&averagepb.AverageResponse{
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

func main() {
	address := "localhost:50051"

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("could not listen on %s: %+v", address, err)
	}

	s := grpc.NewServer()
	averagepb.RegisterAverageServiceServer(s, &server{})

	log.Printf("starting average gRPC server...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("could not serve gRPC server: %+v", err)
	}
}
