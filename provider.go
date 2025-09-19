package main

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Config struct {
	SmtpHost     string
	SmtpPort     int
	SmtpUser     string
	SmtpPassword string

	KeycloakUrl           string
	KeycloakAdminUsername string
	KeycloakAdminPassword string

	HTTPTimeoutSeconds int

	KC *kcClient
}

func Provider() *schema.Provider {
	return &schema.Provider{
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
			"keycloak_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"keycloak_admin_username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"keycloak_admin_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"http_timeout_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"slu_random_password":         resourceSluRandomPassword(),
			"slu_mail_send":               resourceSluMailSend(),
			"slu_keycloak_password_reset": resourceSluKeycloakPasswordReset(),
			"slu_keycloak_org_membership": resourceSluKeycloakOrgMembership(),
		},

		ConfigureContextFunc: func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			var diags diag.Diagnostics

			cfg := &Config{
				SmtpHost:              d.Get("smtp_host").(string),
				SmtpPort:              d.Get("smtp_port").(int),
				SmtpUser:              d.Get("smtp_user").(string),
				SmtpPassword:          d.Get("smtp_password").(string),
				KeycloakUrl:           d.Get("keycloak_url").(string),
				KeycloakAdminUsername: d.Get("keycloak_admin_username").(string),
				KeycloakAdminPassword: d.Get("keycloak_admin_password").(string),
				HTTPTimeoutSeconds:    d.Get("http_timeout_seconds").(int),
			}

			if cfg.KeycloakUrl != "" && cfg.KeycloakAdminUsername != "" && cfg.KeycloakAdminPassword != "" {
				httpClient := &http.Client{Timeout: time.Duration(cfg.HTTPTimeoutSeconds) * time.Second}
				cfg.KC = &kcClient{
					baseURL:   cfg.KeycloakUrl,
					adminUser: cfg.KeycloakAdminUsername,
					adminPass: cfg.KeycloakAdminPassword,
					http:      httpClient,
				}
			}

			return cfg, diags
		},
	}
}
