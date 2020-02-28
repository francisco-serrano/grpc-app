package main

import (
	"github.com/francisco-serrano/grpc-app/decomposition/decompositionpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
}

func (s *server) PrimeNumber(req *decompositionpb.DecompositionRequest, stream decompositionpb.PrimeNumberService_PrimeNumberServer) error {
	log.Printf("received grpc request: %+v\n", req)

	N := req.GetPrimeNumber()
	k := int32(2)

	for N > 1 {
		if N%k == 0 {
			res := &decompositionpb.DecompositionResponse{
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

func main() {
	address := "localhost:50051"

	log.Printf("starting prime number decomposition grpc server at %s\n", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("could not listen at %s: %v", address, err)
	}

	s := grpc.NewServer()
	decompositionpb.RegisterPrimeNumberServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("could not serve grpc server: %v", err)
	}
}
