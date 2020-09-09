package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/sumpb"
	"io"
	"log"
	"math"
	"net"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *sumpb.SumRequest) (*sumpb.SumResponse, error) {
	opr1 := req.GetOperands().GetOperand1()
	opr2 := req.GetOperands().GetOperand2()
	result := opr1 + opr2
	res := &sumpb.SumResponse{
		Result: result,
	}
	return res, nil

}

func (*server) Decomposenumber(req *sumpb.DecomposenumberRequest, stream sumpb.GreetService_DecomposenumberServer) error {
	number := req.GetOperand_Number().GetOperand()
	var k int64 = 2
	for number > 1 {
		if number%k == 0 {
			result := k
			res := &sumpb.DecomposenumberResponse{
				Result: result,
			}
			stream.Send(res)
			number = number / k
			time.Sleep(1000 * time.Millisecond)
		} else {
			k = k + 1
		}
	}
	return nil
}

func (*server) ComputeAverage(stream sumpb.GreetService_ComputeAverageServer) error {
	var num int32 = 0
	count := 0
	var avg float32
	for {
		Request, err := stream.Recv()
		if err == io.EOF {
			avg = float32(num) / float32(count)
			return stream.SendAndClose(&sumpb.ComputeAverageResponse{
				Result: float64(avg),
			})
		}
		if err != nil {
			log.Fatalf("Got error: %v ", err)
		}
		num += Request.GetOperandAvg().GetOperand()
		count++
	}

}

func (*server) Findmax(stream sumpb.GreetService_FindmaxServer) error {
	fmt.Println("Started Bi-direction findmax")
	var max int32 = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("am here : %v", err)
			return err
		}
		fmt.Println("got number", req.GetNum())
		reqNum := req.GetNum()
		if max < reqNum {
			max = reqNum
		}
		fmt.Println("Max", max)
		errr := stream.Send(&sumpb.FindmaxResponse{
			Result: max,
		})

		if errr != nil {
			log.Fatalf("Got error:%v", err)
			return errr
		}
	}

}

func (*server) SquareRoot(ctx context.Context, req *sumpb.SquareRootRequest) (*sumpb.SquareRootResponse, error) {
	fmt.Println("Client invoked SquareRoot")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received the negative number:%v", number),
		)
	}
	return &sumpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Started Server")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	sumpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to server %v", err)
	}
}
