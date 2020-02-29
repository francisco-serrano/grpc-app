package main

import (
	"github.com/francisco-serrano/grpc-app/maximum/maximumpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct {
}

func (s *server) Maximum(stream maximumpb.MaximumService_MaximumServer) error {
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

			res := &maximumpb.MaximumResponse{
				MaxValue: currentMax,
			}

			log.Printf("sending message to client stream %+v\n", res)

			if err := stream.Send(res); err != nil {
				log.Fatalf("error while sending message to client stream: %+v", err)
			}
		}
	}

}

func main() {
	address := "localhost:50051"

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", address, err)
	}

	s := grpc.NewServer()
	maximumpb.RegisterMaximumServiceServer(s, &server{})

	log.Printf("starting maximum gRPC server at %s...\n", address)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("error while attempting to serve gRPC server: %v", err)
	}
}
