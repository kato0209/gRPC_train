package main

import (
	"context"
	"fmt"
	"grpc-lesson/pb"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	certFile := "/home/kato/.local/share/mkcert/rootCA.pem"
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)
	//callListFiles(client)
	callDownload(client)
	//callUpload(client)
	//callUploadAndNotifyProgress(client)
}

func callListFiles(client pb.FileServiceClient) {
	md := metadata.New(map[string]string{"authorization": "Bearer test-token"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	req := &pb.ListFilesRequest{}
	res, err := client.ListFiles(ctx, req)
	if err != nil {
		log.Fatalf("failed to invoke: %v", err)
	}

	fmt.Println(res.GetFilenames())
}

func callDownload(client pb.FileServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.DownloadRequest{
		Filename: "name.txt",
	}
	stream, err := client.Download(ctx, req)
	if err != nil {
		log.Fatalf("failed to invoke: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			resErr, ok := status.FromError(err)
			if ok {
				if resErr.Code() == codes.NotFound {
					log.Fatalf("error code: %v ", "error message: %v", resErr.Code(), resErr.Message())
				} else if resErr.Code() == codes.DeadlineExceeded {
					log.Fatalln("deadline exceeded")
				} else {
					log.Fatalln("unkown grpc error")
				}
			} else {
				log.Fatalf("failed to receive chunk data: %v", err)
			}

		}

		log.Printf("received chunk data: %v", res.GetData())
		log.Printf("received chunk data: %v", string(res.GetData()))
	}
}

func callUpload(client pb.FileServiceClient) {
	filename := "sports.txt"
	path := "/home/kato/Udemy/gRPC/grpc-lesson/storage/" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	defer file.Close()

	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatalf("failed to invoke: %v", err)
	}

	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read chunk data: %v", err)
		}

		req := &pb.UploadRequest{Data: buf[:n]}
		sendErr := stream.Send(req)
		if sendErr != nil {
			log.Fatalf("failed to send chunk data: %v", err)
		}

		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to receive response: %v", err)
	}
	log.Printf("received response: %v", res.GetSize())

}

func callUploadAndNotifyProgress(client pb.FileServiceClient) {
	filename := "sports.txt"
	path := "/home/kato/Udemy/gRPC/grpc-lesson/storage/" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	defer file.Close()

	stream, err := client.UploadAndNotifyProgress(context.Background())
	if err != nil {
		log.Fatalf("failed to invoke: %v", err)
	}

	//request
	buf := make([]byte, 5)
	go func() {
		for {
			n, err := file.Read(buf)
			if n == 0 || err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("failed to read chunk data: %v", err)
			}

			req := &pb.UploadAndNotifyProgressRequest{Data: buf[:n]}
			sendErr := stream.Send(req)
			if sendErr != nil {
				log.Fatalf("failed to send chunk data: %v", err)
			}
			time.Sleep(1 * time.Second)

		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("failed to close send: %v", err)
		}
	}()

	//response
	ch := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("failed to receive chunk data: %v", err)
			}
			log.Printf("received chunk data: %v", res.GetMsg())
		}
		close(ch)
	}()
	<-ch

}
