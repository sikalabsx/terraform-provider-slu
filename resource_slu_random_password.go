package main

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sikalabs/slu/utils/random_utils"
)

func resourceSluRandomPassword() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"result": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},

		Create: resourceSluRandomPasswordCreate,
		Read:   resourceSluRandomPasswordRead,
		Delete: resourceSluRandomPasswordDelete,
	}
}

func resourceSluRandomPasswordCreate(d *schema.ResourceData, m interface{}) error {
	password, err := random_utils.RandomPassword()
	if err != nil {
		return err
	}
	d.SetId(uuid.New().String())
	d.Set("result", password)
	return nil
}

func resourceSluRandomPasswordRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSluRandomPasswordDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
