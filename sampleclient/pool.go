package main

import (
	"context"
	"fmt"
	pb "github.com/karmab/terraform-provider-kcli/kcli-proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

type Kcli struct {
	Url string
}

func main() {
	client := Kcli{Url: "127.0.0.1:50051"}
	conn, err := grpc.Dial(client.Url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	k := pb.NewKcliClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool := pb.Pool{Pool: "images4", Path: "/var/lib/libvirt/images4", Type: "dir"}
	result, err := k.CreatePool(ctx, &pool)
	if err != nil {
		log.Fatalf("could not contact: %v", err)
	}
	fmt.Print(result.Result)
}
