package main

import (
	"context"
	"fmt"
	"grpc-go-course/greet/greetpb"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fiirstName := req.GetGreeting().GetFirstName()
	result := "Hello" + fiirstName
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil

}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fiirstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello" + fiirstName
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreetRequest has been invoked with the reuest")
	result := "Hello"

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//we have received end stream from client
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Here its error:%v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += firstName + "! "
	}

}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone has been invoked with the reuest")
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil

		}
		if err != nil {
			log.Fatalf("We got error:%v", err)
		}
		firstName := request.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "!"
		reterr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if reterr != nil {
			return reterr
		}
	}

}

func (*server) GreetWithDeadLine(ctx context.Context, req *greetpb.GreetWithDeadLineRequest) (*greetpb.GreetWithDeadLineResponse, error) {
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("The client cancelled the request")
			return nil, status.Error(codes.Canceled, "Client canclled the request")
		}
		time.Sleep(1 * time.Second)
	}
	fiirstName := req.GetGreeting().GetFirstName()
	result := "Hello" + fiirstName
	res := &greetpb.GreetWithDeadLineResponse{
		Result: result,
	}

	return res, nil
}

func main() {
	fmt.Println("Server code is here")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to server %v", err)
	}

}
