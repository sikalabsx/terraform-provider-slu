package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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

		CreateContext: resourceSluKeycloakPasswordResetCreate,
		UpdateContext: resourceSluKeycloakPasswordResetCreate,

		ReadContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
			// One-time operation; nothing to read
			return nil
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
			// One-time operation; nothing to delete
			return nil
		},
	}
}

func resourceSluKeycloakPasswordResetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*Config)

	realm := d.Get("realm").(string)
	username := d.Get("username").(string)
	newPassword := d.Get("new_password").(string)

	if err := keycloak_utils.PasswordReset(
		cfg.KeycloakUrl,
		cfg.KeycloakAdminUsername,
		cfg.KeycloakAdminPassword,
		realm,
		username,
		newPassword,
	); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uuid.NewString())
	return nil
}
