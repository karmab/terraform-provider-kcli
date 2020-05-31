package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Kcli struct {
	Url string
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Kcli URL",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"kcli_vm":      resourceVm(),
			"kcli_network": resourceNetwork(),
			"kcli_pool":    resourcePool(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(schema *schema.ResourceData) (interface{}, error) {
	url := schema.Get("url").(string)
	if url == "" {
		url = "127.0.0.1:50051"
	}
	client := Kcli{Url: url}
	return &client, nil
}
