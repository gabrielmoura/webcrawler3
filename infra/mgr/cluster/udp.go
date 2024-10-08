package cluster

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"webCrawler3/config"
	"webCrawler3/infra/pb"
)

func InitCluster(ctx context.Context, conf *config.Config) {
	//grpcServer := grpc.NewServer()
	return

	//pb.RegisterCrawlerServiceServer(grpcServer, &pb.CrawlerServer{})
	//lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.ClusterAddr))
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}
	//
	//go func() {
	//	if err := grpcServer.Serve(lis); err != nil {
	//		log.Fatalf("failed to serve: %v", err)
	//	}
	//}()
}

func RequestQueue(ctx context.Context, conf *config.Config) {
	// inicie grpc client
	conn, err := grpc.Dial(fmt.Sprintf(":%d", conf.ClusterAddr), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewCrawlerServiceClient(conn)

	// Use the client to send requests
	queue, err := client.ProcessCrawl(ctx, &pb.RequestQueue{
		NumberOfProcesses: 10,
	})
	if err != nil {
		log.Fatalf("Error when calling ProcessCrawl: %s", err)
	}
	println(queue.String())
}

func SendProcessed(ctx context.Context, conf *config.Config, result *pb.ProcessedResult) {
	// inicie grpc client
	conn, err := grpc.Dial(fmt.Sprintf(":%d", conf.ClusterAddr), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewCrawlerServiceClient(conn)

	// Use the client to send requests
	_, err = client.ReportProcessedResult(ctx, result)
	if err != nil {
		log.Fatalf("Error when calling ReportProcessedResult: %s", err)
	}
	println("Reported ProcessedResult")
}
