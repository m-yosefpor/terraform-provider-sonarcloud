package sonarcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/m-yosefpor/go-sonarcloud/sonarcloud"
)

// Provider returns a terraform resource provider for SonarCloud.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			// Add any provider-level configs such as tokens, hosts, etc.
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SONARCLOUD_TOKEN", nil),
				Description: "The SonarCloud token used for authentication.",
				Sensitive:   true,
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SonarCloud organization.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sonarcloud_project":             resourceSonarcloudProject(),
			"sonarcloud_qualitygates_select": resourceQualityGatesSelect(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	token := d.Get("token").(string)
	org := d.Get("organization").(string)

	if token == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Token must be provided",
		})
		return nil, diags
	}

	// Create a SonarCloud client
	client := sonarcloud.NewClient(org, token, nil)

	return client, diags
}
