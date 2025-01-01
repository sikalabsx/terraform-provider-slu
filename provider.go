package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Config struct{}

func Provider() *schema.Provider {
	p := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"slu_random_password": resourceSluRandomPassword(),
		},
		Schema: map[string]*schema.Schema{},
	}
	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {

		config := Config{}
		return &config, nil
	}
	return p
}
