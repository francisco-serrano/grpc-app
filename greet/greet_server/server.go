package main

import (
	"context"
	"fmt"
	"github.com/francisco-serrano/grpc-app/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct {
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v", req)

	result := "Hello " + req.GetGreeting().GetFirstName()

	return &greetpb.GreetResponse{
		Result: result,
	}, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", req)

	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: fmt.Sprintf("Hello %s number %s", firstName, strconv.Itoa(i)),
		}

		stream.Send(res)

		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func main() {
	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
