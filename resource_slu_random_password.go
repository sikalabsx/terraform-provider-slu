package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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

		CreateContext: resourceSluRandomPasswordCreate,
		ReadContext:   resourceSluRandomPasswordRead,
		DeleteContext: resourceSluRandomPasswordDelete,
	}
}

func resourceSluRandomPasswordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	password, err := random_utils.RandomPassword()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(uuid.NewString())
	_ = d.Set("result", password)
	return nil
}

func resourceSluRandomPasswordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceSluRandomPasswordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
