package main

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sikalabs/slu/pkg/utils/keycloak_utils"
)

func resourceSluKeycloakPasswordReset() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"new_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},

		Create: resourceSluKeycloakPasswordResetCreate,
		Update: resourceSluKeycloakPasswordResetCreate,
		Read: func(d *schema.ResourceData, m interface{}) error {
			// This resource does not need to be read, as it is a one-time operation.
			return nil
		},
		Delete: func(d *schema.ResourceData, m interface{}) error {
			// This resource does not need to be deleted, as it is a one-time operation.
			return nil
		},
	}
}

func resourceSluKeycloakPasswordResetCreate(d *schema.ResourceData, m interface{}) error {
	err := keycloak_utils.PasswordReset(
		m.(*Config).KeycloakUrl,
		m.(*Config).KeycloakAdminUsername,
		m.(*Config).KeycloakAdminPassword,
		d.Get("realm").(string),
		d.Get("username").(string),
		d.Get("new_password").(string),
	)
	if err != nil {
		return err
	}

	d.SetId(uuid.New().String())
	return nil
}
