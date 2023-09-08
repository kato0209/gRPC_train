package main

import (
	"context"
	"fmt"
	"grpc-lesson/pb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)
	//callListFiles(client)
	callDownload(client)
}

func callListFiles(client pb.FileServiceClient) {
	req := &pb.ListFilesRequest{}
	res, err := client.ListFiles(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to invoke: %v", err)
	}

	fmt.Println(res.GetFilenames())
}

func callDownload(client pb.FileServiceClient) {
	req := &pb.DownloadRequest{
		Filename: "name.txt",
	}
	stream, err := client.Download(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to invoke: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("failed to receive chunk data: %v", err)
		}

		log.Printf("received chunk data: %v", res.GetData())
		log.Printf("received chunk data: %v", string(res.GetData()))
	}
}
