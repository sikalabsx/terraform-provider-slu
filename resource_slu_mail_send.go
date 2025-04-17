package main

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

		Create: resourceSluMailSendCreate,
		Read:   resourceSluMailSendRead,
		Update: resourceSluMailSendUpdate,
		Delete: resourceSluMailSendDelete,
	}
}

func resourceSluMailSendCreate(d *schema.ResourceData, m interface{}) error {
	err := mail_utils.SendSimpleMail(
		m.(*Config).SmtpHost,
		strconv.Itoa(m.(*Config).SmtpPort),
		m.(*Config).SmtpUser,
		m.(*Config).SmtpPassword,
		d.Get("from").(string),
		d.Get("to").(string),
		d.Get("subject").(string),
		d.Get("message").(string),
	)
	if err != nil {
		return err
	}
	d.SetId(uuid.New().String())
	return nil
}

func resourceSluMailSendRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSluMailSendUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSluMailSendDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
