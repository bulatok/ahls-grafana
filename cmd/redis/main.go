package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	clusterOptions := &redis.ClusterOptions{
		Addrs:         []string{"127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003", "127.0.0.1:7004", "127.0.0.1:7005", "127.0.0.1:7006"},
		ReadOnly:      false,
		RouteRandomly: false,
		MaxRedirects:  8,
	}

	// Create a Redis cluster client
	rdb := redis.NewClusterClient(clusterOptions)
	ctx := context.Background()

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		log.Fatal(err)
	}
}
