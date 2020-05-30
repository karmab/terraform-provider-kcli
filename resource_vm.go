package main

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	res, err := config.CreateVm(ctx, vmprofile)
	if err != nil {
		log.Fatalf("could not create vm: %v", err)
	}
	return res
}

func (kcli *Kcli) DeleteVm(vm *pb.Vm) *pb.Result {
	conn, err := grpc.Dial(kcli.Url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	k := pb.NewKcliClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := k.Delete(ctx, vm)
	if err != nil {
		log.Fatalf("could not delete vm: %v", err)
	}
	return res
}

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: createFunc,
		Read:   readFunc,
		Update: updateFunc,
		Delete: deleteFunc,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"image": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"customprofile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"overrides": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ignitionfile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func createFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Kcli)
	vmprofile := pb.Vmprofile{
		Name:         d.Get("name").(string),
		Image:        d.Get("image").(string),
		Overrides:    d.Get("overrides").(string),
		Ignitionfile: d.Get("ignitionfile").(string),
		Profile:      d.Get("profile").(string),
	}

	result := client.CreateVm(&vmprofile)
	if result.Result == "failure" {
		return errors.New(result.Reason)
	}
	d.SetId(vmprofile.Name)
	return nil
}

func readFunc(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func updateFunc(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func deleteFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Kcli)
	vm := pb.Vm{
		Name: d.Get("name").(string),
	}

	result := client.DeleteVm(&vm)
	if result.Result == "failure" {
		return errors.New(result.Reason)
	}
	d.SetId("")
	return nil
}

func providerConfigure(schema *schema.ResourceData) (interface{}, error) {
	url := schema.Get("url").(string)
	if url == "" {
		url = "127.0.0.1:50051"
	}
	client := Kcli{Url: url}
	return &client, nil
}
