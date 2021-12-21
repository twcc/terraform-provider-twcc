package twcc

import (
    "fmt"
    "log"
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIPSecPolicy() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceIPSecPolicyRead,

        Schema: map[string]*schema.Schema{
            "auth_algorithm": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "encapsulation_mode": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "encryption_algorithm": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "transform_protocol": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "lifetime": {
                Type:		schema.TypeInt,
                Computed:	true,
            },

            "pfs": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "name": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },
            "platform": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "project": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
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

// dataSourceIPSecPolicyRead performs the IPSec policy lookup.
func dataSourceIPSecPolicyRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    name := d.Get("name").(string)
    projectID := d.Get("project").(string)
    platform := d.Get("platform").(string)

    resourcePath := fmt.Sprintf("api/v3/%s/ipsec_policies/?project=%s", platform, projectID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to list IPSec policies: %v", err)
    }

    var data []map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    found := false
    for _, ipsec_policy := range data {
        if ipsec_policy["name"] == name {
            if found {
                return fmt.Errorf("There are duplicated IPSec policies with name '%s'", name)
            }
            err = dataSourceIPSecPolicyAttributes(d, ipsec_policy)
            found = true
        }
    }
    if found {
        return err
    }

    return fmt.Errorf("Unable to retrieve IPSec policy %s: %v", name, err)
}

// dataSourceIPSecPolicyAttributes populates the fields of a IPSecPolicy data source.
func dataSourceIPSecPolicyAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    ipsec_policy_id := int(data["id"].(float64))
    log.Printf("[DEBUG] Retrieved twcc_ipsec_policy: %d", ipsec_policy_id)

    d.SetId(fmt.Sprintf("%d", ipsec_policy_id))
    d.Set("auth_algorithm", data["auth_algorithm"])
    d.Set("encryption_algorithm", data["encryption_algorithm"])
    d.Set("encapsulation_mode", data["encapsulation_mode"])
    d.Set("transform_protocol", data["transform_protocol"])
    d.Set("pfs", data["pfs"])
    d.Set("lifetime", data["lifetime"])
    d.Set("user", data["user"])

    return nil
}
