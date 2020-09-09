package main

import (
	"context"
	"fmt"
	"grpc-go-course/blog/blogpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("am your client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("count not connect:%v", err)
	}

	defer conn.Close()
	c := blogpb.NewBlogServiceClient(conn)
	doUnary(c)

	//doReadBlog(c)

	//doUpdateBlog(c)

	//doDeleteBlog(c)

}

func doUnary(c blogpb.BlogServiceClient) {
	fmt.Println("starting to do unary rpc")
	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "Ganesh",
			Title:    "Hello how are you!!!",
			Content:  "Hello Hello Hello!!!!!!",
		},
	}
	res, error := c.CreateBlog(context.Background(), req)

	if error != nil {
		log.Fatalf("Error while calling greet: %v", error)
	}
	log.Printf("Response from greet:%v", res.GetBlog().GetId())

}

func doReadBlog(c blogpb.BlogServiceClient) {
	fmt.Println("starting to do unary rpc")
	req := &blogpb.ReadBlogRequest{
		BlogId: "5f1a81692fcf2cff6de60cf6",
	}
	res, error := c.ReadBlog(context.Background(), req)
	if error != nil {
		log.Fatalf("We got error while reading")
	}

	log.Printf("Author :%v , title :%v , content %v", res.GetBlog().GetAuthorId(), res.GetBlog().GetContent(), res.GetBlog().GetTitle())

}

func doUpdateBlog(c blogpb.BlogServiceClient) {
	fmt.Println("starting to do unary rpc")
	req := &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       "5f1a81692fcf2cff6de60cf6",
			AuthorId: "Hulkat",
			Title:    "How is world",
			Content:  "This is the content of world",
		},
	}
	res, error := c.UpdateBlog(context.Background(), req)
	if error != nil {
		log.Fatalf("We got error while reading")
	}

	log.Printf("Author :%v , title :%v , content %v", res.GetBlog().GetAuthorId(), res.GetBlog().GetContent(), res.GetBlog().GetTitle())

}

func doDeleteBlog(c blogpb.BlogServiceClient) {
	req := &blogpb.DeleteBlogRequest{
		BlogId: "5f1a81692fcf2cff6de60cf6",
	}
	res, err := c.DeleteBlog(context.Background(), req)
	if err != nil {
		fmt.Printf("Not able to delete data :%v", err)
	}

	fmt.Printf("Blog was deleted: %v \n", res)

}
