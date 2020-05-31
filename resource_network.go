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

func (kcli *Kcli) CreateNetwork(network *pb.Network) *pb.Result {
	conn, err := grpc.Dial(kcli.Url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	k := pb.NewKcliClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := k.CreateNetwork(ctx, network)
	if err != nil {
		log.Fatalf("could not create vm: %v", err)
	}
	return res
}

func (kcli *Kcli) DeleteNetwork(network *pb.Network) *pb.Result {
	conn, err := grpc.Dial(kcli.Url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	k := pb.NewKcliClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := k.DeleteNetwork(ctx, network)
	if err != nil {
		log.Fatalf("could not delete vm: %v", err)
	}
	return res
}

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Create: NetworkcreateFunc,
		Read:   NetworkreadFunc,
		Update: NetworkupdateFunc,
		Delete: NetworkdeleteFunc,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"dhcp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"overrides": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func NetworkcreateFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Kcli)
	network := pb.Network{
		Network:   d.Get("name").(string),
		Cidr:      d.Get("cidr").(string),
		Dhcp:      d.Get("dhcp").(string),
		Domain:    d.Get("domain").(string),
		Overrides: d.Get("overrides").(string),
	}
	result := client.CreateNetwork(&network)
	if result.Result == "failure" {
		return errors.New(result.Reason)
	}
	d.SetId(network.Network)
	return nil
}

func NetworkreadFunc(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func NetworkupdateFunc(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func NetworkdeleteFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Kcli)
	network := pb.Network{
		Network: d.Get("name").(string),
	}

	result := client.DeleteNetwork(&network)
	if result.Result == "failure" {
		return errors.New(result.Reason)
	}
	d.SetId("")
	return nil
}
