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

func (kcli *Kcli) CreatePool(pool *pb.Pool) *pb.Result {
	conn, err := grpc.Dial(kcli.Url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	k := pb.NewKcliClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := k.CreatePool(ctx, pool)
	if err != nil {
		log.Fatalf("could not create pool: %v", err)
	}
	return res
}

func (kcli *Kcli) DeletePool(pool *pb.Pool) *pb.Result {
	conn, err := grpc.Dial(kcli.Url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	k := pb.NewKcliClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := k.DeletePool(ctx, pool)
	if err != nil {
		log.Fatalf("could not delete pool: %v", err)
	}
	return res
}

func resourcePool() *schema.Resource {
	return &schema.Resource{
		Create: PoolcreateFunc,
		Read:   PoolreadFunc,
		Update: PoolupdateFunc,
		Delete: PooldeleteFunc,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"thin": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func PoolcreateFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Kcli)
	pool := pb.Pool{
		Pool:     d.Get("name").(string),
		Type:     d.Get("type").(string),
		Path:     d.Get("path").(string),
		Thinpool: d.Get("thin").(bool),
	}
	_type := d.Get("type").(string)
	if _type == "" {
		pool.Type = "dir"
	}
	result := client.CreatePool(&pool)
	if result.Result == "failure" {
		return errors.New(result.Reason)
	}
	d.SetId(pool.Pool)
	return nil
}

func PoolreadFunc(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func PoolupdateFunc(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func PooldeleteFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Kcli)
	pool := pb.Pool{
		Pool: d.Get("name").(string),
	}

	result := client.DeletePool(&pool)
	if result.Result == "failure" {
		return errors.New(result.Reason)
	}
	d.SetId("")
	return nil
}
