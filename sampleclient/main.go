package main

import (
	"context"
	pb "github.com/karmab/terraform-provider-kcli/kcli-proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

type Kcli struct {
	Url string
}

func (kcli *Kcli) CreateVm(vmprofile *pb.Vmprofile) *pb.Result {
	conn, err := grpc.Dial(kcli.Url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	config := pb.NewKconfigClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// res, err := config.CreateVm(ctx, &pb.Vmprofile{Name: vmprofile.Name, Image: vmprofile.Image})
	res, err := config.CreateVm(ctx, vmprofile)
	if err != nil {
		log.Fatalf("could not contact: %v", err)
	}
	return res
}

func main() {
	client := Kcli{Url: "127.0.0.1:50051"}
	name := "federer"
	vmprofile := pb.Vmprofile{
		Name:    name,
		Profile: "ubuntu1910",
	}
	result := client.CreateVm(&vmprofile)
	if result.Result == "success" {
		log.Printf("VM %s was created", result.Vm)
	} else {
		log.Printf("VM %s not created because %s", name, result.Reason)
	}
}
