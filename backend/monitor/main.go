package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "tanken/data-fetcher" // 更新为实际路径
)

const (
	redisAddr       = "geo-postid-redis:6379"
	grpcServerAddr  = "data-fetcher:50051"
	memoryThreshold = 60.0
	checkInterval   = 10 * time.Second
)

func main() {
	geo_postid_rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	conn, err := grpc.Dial(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewDataFetcherClient(conn)

	// 开始监控循环
	monitorRedisMemoryUsage(geo_postid_rdb, client)
}

func monitorRedisMemoryUsage(rdb *redis.Client, client pb.DataFetcherClient) {
	ctx := context.Background()
	for {
		// 获取 Redis 内存信息
		usage, err := getMemoryUsage(ctx, rdb)
		if err != nil {
			log.Printf("Error fetching Redis memory usage: %v", err)
			continue
		}

		fmt.Printf("Current Redis memory usage: %.2f%%\n", usage)

		// 检查是否超过阈值
		if usage > memoryThreshold {
			fmt.Println("Memory threshold exceeded. Notifying data-fetcher...")
			notifyDataFetcher(ctx, client, usage)
		}

		time.Sleep(checkInterval)
	}
}

func getMemoryUsage(ctx context.Context, rdb *redis.Client) (float64, error) {
	result, err := rdb.Info(ctx, "memory").Result()
	if err != nil {
		return 0, err
	}
	// 解析结果，获取 used_memory 和 maxmemory
	var usedMemory, maxMemory float64
	fmt.Sscanf(result, "used_memory:%f\r\nmaxmemory:%f\r\n", &usedMemory, &maxMemory)
	if maxMemory == 0 {
		return 0, fmt.Errorf("maxmemory not defined")
	}
	return (usedMemory / maxMemory) * 100, nil
}

func notifyDataFetcher(ctx context.Context, client pb.DataFetcherClient, usage float64) {
	req := &pb.WriteBackCacheRequest{
		Usage: float32(usage),
	}
	if _, err := client.WriteBackCache(ctx, req); err != nil {
		log.Printf("Failed to notify data-fetcher: %v", err)
	} else {
		log.Println("Data-fetcher notified successfully.")
	}
}
