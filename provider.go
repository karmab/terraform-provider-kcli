package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

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
			"kcli_vm": resourceServer(),
		},
		ConfigureFunc: providerConfigure,
	}
}
