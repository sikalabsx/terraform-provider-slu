package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSluKeycloakOrgMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSluKeycloakOrgMembershipCreate,
		ReadContext:   resourceSluKeycloakOrgMembershipRead,
		DeleteContext: resourceSluKeycloakOrgMembershipDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSluKeycloakOrgMembershipImport,
		},

		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"exists": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func kcFromMeta(m interface{}) (*kcClient, diag.Diagnostics) {
	cfg := m.(*Config)
	if cfg.KC == nil {
		return nil, diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Keycloak not configured",
				Detail:   "keycloak_url, keycloak_admin_username and keycloak_admin_password must be set in the provider to use slu_keycloak_org_membership.",
			},
		}
	}
	return cfg.KC, nil
}

func resourceSluKeycloakOrgMembershipCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	kc, diags := kcFromMeta(m)
	if diags.HasError() {
		return diags
	}

	realm := d.Get("realm_id").(string)
	org := d.Get("org_id").(string)
	username := d.Get("username").(string)

	userID := d.Get("user_id").(string)
	if userID == "" {
		id, err := kc.ResolveUserID(ctx, realm, username)
		if err != nil {
			return diag.FromErr(fmt.Errorf("resolve user id for %q: %w", username, err))
		}
		userID = id
	}

	if err := kc.AddMember(ctx, realm, org, userID); err != nil {
		return diag.FromErr(fmt.Errorf("add member %q to org %q: %w", userID, org, err))
	}

	// Stable resource ID; if you prefer deterministic over uuid, you can use realm/org/userID
	d.SetId(fmt.Sprintf("%s/%s/%s", realm, org, userID))
	_ = d.Set("user_id", userID)
	_ = d.Set("exists", true)

	return resourceSluKeycloakOrgMembershipRead(ctx, d, m)
}

func resourceSluKeycloakOrgMembershipRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	kc, diags := kcFromMeta(m)
	if diags.HasError() {
		return diags
	}

	realm := d.Get("realm_id").(string)
	org := d.Get("org_id").(string)
	userID := d.Get("user_id").(string)

	if userID == "" {
		// In case of import with username only
		if username, ok := d.GetOk("username"); ok {
			if resolved, err := kc.ResolveUserID(ctx, realm, username.(string)); err == nil && resolved != "" {
				userID = resolved
				_ = d.Set("user_id", userID)
			}
		}
	}

	ok, err := kc.CheckMember(ctx, realm, org, userID)
	if err != nil {
		if apiErr, ok2 := err.(apiError); ok2 && apiErr.Status == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("check membership %s/%s/%s: %w", realm, org, userID, err))
	}

	_ = d.Set("exists", ok)
	if !ok {
		d.SetId("")
	}
	return nil
}

func resourceSluKeycloakOrgMembershipDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	kc, diags := kcFromMeta(m)
	if diags.HasError() {
		return diags
	}

	realm := d.Get("realm_id").(string)
	org := d.Get("org_id").(string)
	userID := d.Get("user_id").(string)

	if err := kc.RemoveMember(ctx, realm, org, userID); err != nil {
		if apiErr, ok := err.(apiError); ok && apiErr.Status == http.StatusNotFound {
			// already gone
		} else {
			return diag.FromErr(fmt.Errorf("remove membership %s/%s/%s: %w", realm, org, userID, err))
		}
	}

	d.SetId("")
	return nil
}

// Import ID formats supported:
//
//	<realm>/<org>/<userId>  OR  <realm>/<org>/<username>
func resourceSluKeycloakOrgMembershipImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	kc, diags := kcFromMeta(m)
	if diags.HasError() {
		return nil, fmt.Errorf(diags[0].Detail)
	}

	var realm, org, user string
	n, _ := fmt.Sscanf(d.Id(), "%s/%s/%s", &realm, &org, &user)
	if n != 3 {
		return nil, fmt.Errorf("import id must be <realm>/<org>/<userId|username>")
	}

	_ = d.Set("realm_id", realm)
	_ = d.Set("org_id", org)
	_ = d.Set("username", user)

	if resolved, err := kc.ResolveUserID(ctx, realm, user); err == nil && resolved != "" {
		_ = d.Set("user_id", resolved)
		d.SetId(fmt.Sprintf("%s/%s/%s", realm, org, resolved))
		return []*schema.ResourceData{d}, nil
	}

	// fallback: treat provided user as the id
	_ = d.Set("user_id", user)
	d.SetId(fmt.Sprintf("%s/%s/%s", realm, org, user))
	return []*schema.ResourceData{d}, nil
}
