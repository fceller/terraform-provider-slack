package slack

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
	"log"
)

func dataSourceUserGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSlackUserGroupRead,

		Schema: map[string]*schema.Schema{
			"usergroup_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"handle": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"auto_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSlackUserGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Team).client

	usergroupId := d.Get("usergroup_id").(string)

	log.Printf("[DEBUG] Reading usergroup: %s", usergroupId)
	groups, err := client.GetUserGroupsContext(ctx, func(params *slack.GetUserGroupsParams) {
		params.IncludeUsers = false
		params.IncludeCount = false
		params.IncludeDisabled = true
	})

	if err != nil {
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("provicer cannot find a usergroup (%s)", usergroupId),
				Detail:   err.Error(),
			},
		}
	}

	for _, group := range groups {
		if group.ID == usergroupId {
			d.SetId(group.ID)
			_ = d.Set("handle", group.Handle)
			_ = d.Set("name", group.Name)
			_ = d.Set("description", group.Description)
			_ = d.Set("auto_type", group.AutoType)
			_ = d.Set("team_id", group.TeamID)
			return nil
		}
	}

	return diag.Diagnostics{
		{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("provicer cannot find a usergroup (%s)", usergroupId),
			Detail:   fmt.Sprintf("a usergroup (%s) is not found in available usergroups that this token can view", usergroupId),
		},
	}
}
