package twcc

import (
    "fmt"
    "log"
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAutoScalingPolicy() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceAutoScalingPolicyRead,

        Schema: map[string]*schema.Schema{
            "name": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "project": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "platform": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "meter_name": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "description": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "scale_max_size": {
                Type:		schema.TypeInt,
                Computed:	true,
            },

            "scaledown_threshold": {
                Type:		schema.TypeInt,
                Computed:	true,
            },

            "scaleup_threshold": {
                Type:		schema.TypeInt,
                Computed:	true,
            },

            "user": {
                Type:		schema.TypeMap,
                Computed:	true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
        },
    }
}

// dataSourceAutoScalingPolicyRead performs the auto scaling policy lookup.
func dataSourceAutoScalingPolicyRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    name := d.Get("name").(string)
    projectID := d.Get("project").(string)
    platform := d.Get("platform").(string)

    resourcePath := fmt.Sprintf("api/v3/%s/auto_scaling_policies/?project=%s&name=%s", platform, projectID, name)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to list auto scaling policies: %v", err)
    }

    var data []map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    for _, auto_scaling_policy := range data {
        if auto_scaling_policy["name"] == name {
            return dataSourceAutoScalingPolicyAttributes(d, auto_scaling_policy)
        }
    }

    return fmt.Errorf("Unable to retrieve auto scaling policy %s: %v", name, err)
}

// dataSourceAutoScalingPolicyAttributes populates the fields of a auto scaling policy data source.
func dataSourceAutoScalingPolicyAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    auto_scaling_policy_id := int(data["id"].(float64))
    log.Printf("[DEBUG] Retrieved twcc_auto_scaling_policy: %d", auto_scaling_policy_id)

    d.SetId(fmt.Sprintf("%d", auto_scaling_policy_id))
    d.Set("meter_name", data["meter_name"])
    d.Set("description", data["description"])
    d.Set("scale_max_size", data["scale_max_size"])
    d.Set("scaledown_threshold", data["scaledown_threshold"])
    d.Set("scaleup_threshold", data["scaleup_threshold"])
    d.Set("user", data["user"])

    return nil
}
