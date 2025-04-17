package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Config struct {
	SmtpHost     string
	SmtpPort     int
	SmtpUser     string
	SmtpPassword string
}

func Provider() *schema.Provider {
	p := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"slu_random_password": resourceSluRandomPassword(),
			"slu_mail_send":       resourceSluMailSend(),
		},
		Schema: map[string]*schema.Schema{
			"smtp_host": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"smtp_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"smtp_user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"smtp_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {

		config := Config{
			SmtpHost:     d.Get("smtp_host").(string),
			SmtpPort:     d.Get("smtp_port").(int),
			SmtpUser:     d.Get("smtp_user").(string),
			SmtpPassword: d.Get("smtp_password").(string),
		}
		return &config, nil
	}
	return p
}
