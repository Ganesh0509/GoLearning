package main

import (
	"fmt"
	"grpc-go-course/"
)

func main() {
 fmt.Println("Starting client!! Make sure server is availble!!!!")
 stream1.pb.CreateBlog(ctx context.Context, opts ...grpc.CallOption) (SimpletestService_CreateBlogClient, error)
}
