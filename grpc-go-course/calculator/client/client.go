package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/sumpb"
	"io"
	"log"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Am your client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("count not connect:%v", err)
	}
	defer conn.Close()
	c := sumpb.NewGreetServiceClient(conn)

	fmt.Println("client end")
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBidirectionStreamin(c)
	doErrorUnary(c)
}

func doUnary(c sumpb.GreetServiceClient) {
	fmt.Println("starting to do unary rpc")
	req := &sumpb.SumRequest{
		Operands: &sumpb.Operands{
			Operand1: 5,
			Operand2: 6,
		},
	}
	res, error := c.Sum(context.Background(), req)
	if error != nil {
		log.Fatal("Error while calling greet %v", error)
	}
	log.Println("Response from greet:%v", res.Result)
}

func doServerStreaming(c sumpb.GreetServiceClient) {
	fmt.Println("Started streaming")
	req := &sumpb.DecomposenumberRequest{
		Operand_Number: &sumpb.Operand_Number{
			Operand: 210,
		},
	}

	stm, err := c.Decomposenumber(context.Background(), req)
	if err != nil {
		log.Fatalf("error", err)
	}
	for {
		msg, err := stm.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error", err)
		}
		log.Printf("Message", msg)

	}
}

func doClientStreaming(c sumpb.GreetServiceClient) {
	request := []*sumpb.ComputeAverageRequest{
		&sumpb.ComputeAverageRequest{
			OperandAvg: &sumpb.OperandAvg{
				Operand: 1,
			},
		},
		&sumpb.ComputeAverageRequest{
			OperandAvg: &sumpb.OperandAvg{
				Operand: 2,
			},
		},
		&sumpb.ComputeAverageRequest{
			OperandAvg: &sumpb.OperandAvg{
				Operand: 3,
			},
		},
		&sumpb.ComputeAverageRequest{
			OperandAvg: &sumpb.OperandAvg{
				Operand: 4,
			},
		},
	}
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Got errr:%v", err)
	}

	for _, req := range request {
		stream.Send(req)
	}
	respo, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("We got error : %v", err)
	}
	fmt.Printf("response :%v", respo.GetResult())

}

func doBidirectionStreamin(c sumpb.GreetServiceClient) {

	request := []*sumpb.FindmaxRequest{
		&sumpb.FindmaxRequest{
			Num: 1,
		},
		&sumpb.FindmaxRequest{
			Num: 5,
		},
		&sumpb.FindmaxRequest{
			Num: 3,
		},
		&sumpb.FindmaxRequest{
			Num: 6,
		},
		&sumpb.FindmaxRequest{
			Num: 2,
		},
		&sumpb.FindmaxRequest{
			Num: 20,
		},
	}

	stream, err := c.Findmax(context.Background())
	if err != nil {
		log.Fatalf("Got error while invoking Server:%v", err)
		return
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range request {
			err1 := stream.Send(req)
			if err1 != nil {
				log.Fatalf("Got error when sending request:%v", err1)
				return
			}
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
			}
			if err != nil {
				log.Fatalf("error occurent :%v", err)
				close(waitc)
			}
			fmt.Printf("Response:%v\n", res.GetResult())
		}

	}()

	<-waitc

}

func doErrorUnary(c sumpb.GreetServiceClient) {
	req := &sumpb.SquareRootRequest{
		Number: -10,
	}

	response, err := c.SquareRoot(context.Background(), req)
	if err != nil {
		respErr, ok := status.FromError(err)

		if ok {
			fmt.Println("User error", respErr.Message())
			fmt.Println("User error", respErr.Err())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent invalid data")
			}
		} else {
			log.Fatalf("Big error: %v", respErr)
			return
		}
	}

	fmt.Println(" We got success", response.GetNumberRoot())
}
