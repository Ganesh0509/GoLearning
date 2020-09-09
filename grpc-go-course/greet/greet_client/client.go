package main

import (
	"context"
	"fmt"
	"grpc-go-course/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("am your client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("count not connect:%v", err)
	}

	defer conn.Close()
	c := greetpb.NewGreetServiceClient(conn)
	//doUnary(c)
	//doServerStreaming(c)
	//doclientStreaming(c)
	//doBiDiStreaming(c)
	doUnaryWithDeadline(c, 5*time.Second)
	doUnaryWithDeadline(c, 1*time.Second)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do unary rpc")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName:  "Ganesh",
			SecondName: "Dixith",
		},
	}
	res, error := c.Greet(context.Background(), req)

	if error != nil {
		log.Fatal("Error while calling greet %v", error)
	}
	log.Println("Response from greet:%v", res.Result)

}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Server streaming RPC")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName:  "Ganesh",
			SecondName: "Dixith",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("error", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error", err)
		}
		log.Printf("Message", msg)
	}

}

func doclientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("client streaming RPC")
	request := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName:  "Gani",
				SecondName: "Test",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName:  "Gani",
				SecondName: "Test",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName:  "Gani",
				SecondName: "Test",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Got errr:%v", err)
	}
	for _, req := range request {
		fmt.Printf("semding request : %v", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	respo, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("We got error : %v", err)
	}
	fmt.Printf("response :%v", respo.GetResult())
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Started doBiDiStreaming")
	request := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName:  "Gani",
				SecondName: "Test",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName:  "Gani",
				SecondName: "Test",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName:  "Gani",
				SecondName: "Test",
			},
		},
	}

	//We create stream by invoking client
	strm, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	waitc := make(chan struct{})

	//we send bunch of messages
	go func() {
		for _, req := range request {
			fmt.Printf("semding request : %v\n", req)
			err := strm.Send(req)
			if err != nil {
				fmt.Println("Got error", err)
			}
			fmt.Printf("semding request2 : %v\n", req)
			time.Sleep(100 * time.Millisecond)
		}
		strm.CloseSend()
	}()

	//We receive bunch of messages
	go func() {
		for {
			response, err := strm.Recv()
			if err == io.EOF {
				close(waitc)
			}
			if err != nil {
				log.Fatalf("Got error: %v", err)
				close(waitc)
			}
			fmt.Printf("Response\n", response)
		}
	}()
	//block entill everything is done
	<-waitc

}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("starting to do unary rpc")
	req := &greetpb.GreetWithDeadLineRequest{
		Greeting: &greetpb.Greeting{
			FirstName:  "Ganesh",
			SecondName: "Dixith",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, error := c.GreetWithDeadLine(ctx, req)

	if error != nil {
		statusErr, ok := status.FromError(error)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout exeeded ")
			} else {
				fmt.Println("unexpected error", statusErr)
			}

		} else {
			log.Fatalf("Error while calling greet :%v", error)
		}
		return
	}
	log.Printf("Response from greet:%v", res.GetResult())

}
