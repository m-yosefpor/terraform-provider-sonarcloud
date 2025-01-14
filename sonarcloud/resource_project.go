package sonarcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/m-yosefpor/go-sonarcloud/sonarcloud"
	"github.com/m-yosefpor/go-sonarcloud/sonarcloud/paging"
	"github.com/m-yosefpor/go-sonarcloud/sonarcloud/projects"
)

func resourceSonarcloudProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSonarcloudProjectCreate,
		ReadContext:   resourceSonarcloudProjectRead,
		UpdateContext: resourceSonarcloudProjectUpdate,
		DeleteContext: resourceSonarcloudProjectDelete,

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The key of the organization under which the project is created.",
				ForceNew:    true,
			},
			"project_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique key of the project.",
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the project.",
			},
			"branch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SCM Branch of the project. If provided, final key is project_key:branch.",
				ForceNew:    true,
			},
			"visibility": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Project visibility: 'public' or 'private'. If not set, default org visibility is used.",
			},
			"new_code_definition_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the new code definition. One of 'previous_version', 'days', 'date', 'version'.",
			},
			"new_code_definition_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Value for the new code definition, depends on the chosen type.",
			},
		},
	}
}

func resourceSonarcloudProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*sonarcloud.Client)

	org := d.Get("organization").(string)
	key := d.Get("project_key").(string)
	name := d.Get("name").(string)
	branch := d.Get("branch").(string)
	visibility := d.Get("visibility").(string)
	ncdType := d.Get("new_code_definition_type").(string)
	ncdValue := d.Get("new_code_definition_value").(string)

	// Check if the project already exists
	searchReq := projects.SearchRequest{
		Organization: org,
		Projects:     key,
	}
	searchResp, err := client.Projects.Search(searchReq, paging.Params{})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error checking if sonarcloud project exists: %w", err))
	}

	if len(searchResp.Components) > 0 {
		// Project already exists, set the ID and return
		d.SetId(searchResp.Components[0].Key)
		return resourceSonarcloudProjectRead(ctx, d, m)
	}

	req := projects.CreateRequest{
		Organization:           org,
		Project:                key,
		Name:                   name,
		Branch:                 branch,
		Visibility:             visibility,
		NewCodeDefinitionType:  ncdType,
		NewCodeDefinitionValue: ncdValue,
	}

	resp, err := client.Projects.Create(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating sonarcloud project: %w", err))
	}

	d.SetId(resp.Project.Key)
	return resourceSonarcloudProjectRead(ctx, d, m)
}

func resourceSonarcloudProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*sonarcloud.Client)

	org := d.Get("organization").(string)
	projectKey := d.Id()

	searchReq := projects.SearchRequest{
		Organization: org,
		Projects:     projectKey,
	}
	searchResp, err := client.Projects.Search(searchReq, paging.Params{})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading sonarcloud project: %w", err))
	}

	if len(searchResp.Components) == 0 {
		// Project not found
		d.SetId("")
		return nil
	}

	comp := searchResp.Components[0]

	// Set data back to state
	if err := d.Set("name", comp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("organization", comp.Organization); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_key", comp.Key); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("visibility", comp.Visibility); err != nil {
		return diag.FromErr(err)
	}

	// Note: Branch, new_code_definition_* are not retrieved by Search endpoint directly.
	// If you need them, you'll have to store them in TF state or fetch from another endpoint.
	// For now, we rely on what was set and don't overwrite.

	return nil
}

func resourceSonarcloudProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*sonarcloud.Client)

	if d.HasChange("visibility") {
		projectKey := d.Id()
		visibility := d.Get("visibility").(string)

		updateReq := projects.UpdateVisibilityRequest{
			Project:    projectKey,
			Visibility: visibility,
		}

		if err := client.Projects.UpdateVisibility(updateReq); err != nil {
			return diag.FromErr(fmt.Errorf("error updating sonarcloud project visibility: %w", err))
		}
	}

	// Name, NewCodeDefinition, and Branch updates are not shown as supported by this snippet.
	// If the SonarCloud API supports updating these, call the appropriate endpoint here.
	// Otherwise, if changes are detected to fields that cannot be updated, you may need to force recreate.

	return resourceSonarcloudProjectRead(ctx, d, m)
}

func resourceSonarcloudProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*sonarcloud.Client)
	projectKey := d.Id()

	delReq := projects.DeleteRequest{
		Project: projectKey,
	}

	if err := client.Projects.Delete(delReq); err != nil {
		return diag.FromErr(fmt.Errorf("error deleting sonarcloud project: %w", err))
	}

	d.SetId("")
	return nil
}
