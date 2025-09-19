package main

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sikalabs/slu/utils/mail_utils"
)

func resourceSluMailSend() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"from": {
				Type:     schema.TypeString,
				Required: true,
			},
			"message": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subject": {
				Type:     schema.TypeString,
				Required: true,
			},
			"to": {
				Type:     schema.TypeString,
				Required: true,
			},
		},

		CreateContext: resourceSluMailSendCreate,
		ReadContext:   resourceSluMailSendRead,
		UpdateContext: resourceSluMailSendUpdate,
		DeleteContext: resourceSluMailSendDelete,
	}
}

func resourceSluMailSendCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*Config)

	err := mail_utils.SendSimpleMail(
		cfg.SmtpHost,
		strconv.Itoa(cfg.SmtpPort),
		cfg.SmtpUser,
		cfg.SmtpPassword,
		d.Get("from").(string),
		d.Get("to").(string),
		d.Get("subject").(string),
		d.Get("message").(string),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uuid.New().String())
	return nil
}

func resourceSluMailSendRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceSluMailSendUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceSluMailSendDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
