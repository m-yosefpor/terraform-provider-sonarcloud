package sonarcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/m-yosefpor/go-sonarcloud/sonarcloud"
	"github.com/m-yosefpor/go-sonarcloud/sonarcloud/qualitygates"
)

func resourceQualityGatesSelect() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceQualityGatesSelectCreate,
		ReadContext:   resourceQualityGatesSelectRead,
		DeleteContext: resourceQualityGatesSelectDelete,

		Schema: map[string]*schema.Schema{
			"project_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The key of the project.",
			},
			"quality_gate_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the quality gate to select.",
			},
		},
	}
}

func resourceQualityGatesSelectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*sonarcloud.Client)

	projectKey := d.Get("project_key").(string)
	qualityGateID := d.Get("quality_gate_id").(int)

	req := qualitygates.SelectRequest{
		ProjectKey:    projectKey,
		QualityGateID: qualityGateID,
	}

	if err := client.QualityGates.Select(req); err != nil {
		return diag.FromErr(fmt.Errorf("error selecting quality gate: %w", err))
	}

	d.SetId(fmt.Sprintf("%s-%d", projectKey, qualityGateID))
	return resourceQualityGatesSelectRead(ctx, d, m)
}

func resourceQualityGatesSelectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*sonarcloud.Client)

	projectKey := d.Get("project_key").(string)

	req := qualitygates.GetByProjectRequest{
		ProjectKey: projectKey,
	}

	resp, err := client.QualityGates.GetByProject(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading quality gate: %w", err))
	}

	if err := d.Set("quality_gate_id", resp.QualityGate.ID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceQualityGatesSelectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*sonarcloud.Client)

	projectKey := d.Get("project_key").(string)

	req := qualitygates.DeselectRequest{
		ProjectKey: projectKey,
	}

	if err := client.QualityGates.Deselect(req); err != nil {
		return diag.FromErr(fmt.Errorf("error deselecting quality gate: %w", err))
	}

	d.SetId("")
	return nil
}
