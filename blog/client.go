package main

import (
	"context"
	"github.com/francisco-serrano/grpc-app/blog/blogpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	log.Println("Bloc Client")

	opts := grpc.WithInsecure()

	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()

	c := blogpb.NewBlogServiceClient(conn)

	// Create Blog
	blog := &blogpb.Blog{
		AuthorId: "Francisco",
		Title:    "My First Blog",
		Content:  "Content of the First Blog",
	}

	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	log.Printf("blog has been created: %v\n", res)
	blogId := res.GetBlog().GetId()

	// Reading Blog
	log.Println("Reading the Blog")

	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "5e5b4cc3e389cb3072051bd7",
	})
	if err2 != nil {
		log.Printf("error happenned while reading: %v\n", err2)
	}

	readBlogReq := &blogpb.ReadBlogRequest{
		BlogId: blogId,
	}

	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		log.Printf("error happenned while reading: %v\n", readBlogErr)
	}

	log.Printf("blog was read: %v\n", readBlogRes)

	// Updating Blog
	newBlog := &blogpb.Blog{
		Id:       blogId,
		AuthorId: "Changed Author",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}

	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: newBlog,
	})

	if updateErr != nil {
		log.Printf("error while updating %v\n", updateErr)
	}

	log.Printf("blog was updated: %v", updateRes)

	// Deleting Blog
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: blogId,
	})
	if deleteErr != nil {
		log.Printf("error while deleting: %v\n", deleteErr)
	}

	log.Printf("Blog was deleted: %v \n", deleteRes)

	// Listing Blogs
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("something happenned: %v", err)
		}

		log.Printf("blog was read: %v\n", res.GetBlog())
	}
}
