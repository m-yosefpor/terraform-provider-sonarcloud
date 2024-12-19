package sonarcloud

import (
	"context"
	"fmt"
	"strconv"

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
				ForceNew:    true,
			},
			"quality_gate_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the quality gate to select.",
				ForceNew:    true,
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of organization.",
				ForceNew:    true,
			},
		},
	}
}

func resourceQualityGatesSelectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*sonarcloud.Client)

	projectKey := d.Get("project_key").(string)
	qualityGateID := d.Get("quality_gate_id").(int)
	organization := d.Get("organization").(string)

	req := qualitygates.SelectRequest{
		ProjectKey:   projectKey,
		Organization: organization,
		GateId:       strconv.Itoa(qualityGateID),
	}

	if err := client.Qualitygates.Select(req); err != nil {
		return diag.FromErr(fmt.Errorf("error selecting quality gate: %w", err))
	}

	d.SetId(fmt.Sprintf("%s-%d", projectKey, qualityGateID))
	return resourceQualityGatesSelectRead(ctx, d, m)
}

func resourceQualityGatesSelectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*sonarcloud.Client)

	projectKey := d.Get("project_key").(string)
	organization := d.Get("organization").(string)

	req := qualitygates.GetByProjectRequest{
		Project:      projectKey,
		Organization: organization,
	}

	resp, err := client.Qualitygates.GetByProject(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading quality gate: %w", err))
	}

	if err := d.Set("quality_gate_id", resp.QualityGate.Id); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceQualityGatesSelectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*sonarcloud.Client)

	projectKey := d.Get("project_key").(string)
	organization := d.Get("organization").(string)

	req := qualitygates.DeselectRequest{
		ProjectKey:   projectKey,
		Organization: organization,
	}

	if err := client.Qualitygates.Deselect(req); err != nil {
		return diag.FromErr(fmt.Errorf("error deselecting quality gate: %w", err))
	}

	d.SetId("")
	return nil
}
