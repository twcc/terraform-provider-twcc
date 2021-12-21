package twcc

import (
    "fmt"
    "log"
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIKEPolicy() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceIKEPolicyRead,

        Schema: map[string]*schema.Schema{
            "auth_algorithm": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "encryption_algorithm": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "ike_version": {
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

// dataSourceIKEPolicyRead performs the IKE policy lookup.
func dataSourceIKEPolicyRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    name := d.Get("name").(string)
    projectID := d.Get("project").(string)
    platform := d.Get("platform").(string)

    resourcePath := fmt.Sprintf("api/v3/%s/ike_policies/?project=%s", platform, projectID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to list IKE policies: %v", err)
    }

    var data []map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    found := false
    for _, ike_policy := range data {
        if ike_policy["name"] == name {
            if found {
                return fmt.Errorf("There are duplicated IKE policies with name '%s'", name)
            }
            err = dataSourceIKEPolicyAttributes(d, ike_policy)
            found = true
        }
    }
    if found {
        return err
    }

    return fmt.Errorf("Unable to retrieve IKE policy %s: %v", name, err)
}

// dataSourceIKEPolicyAttributes populates the fields of a IKEPolicy data source.
func dataSourceIKEPolicyAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    ike_policy_id := int(data["id"].(float64))
    log.Printf("[DEBUG] Retrieved twcc_ike_policy: %d", ike_policy_id)

    d.SetId(fmt.Sprintf("%d", ike_policy_id))
    d.Set("auth_algorithm", data["auth_algorithm"])
    d.Set("ike_version", data["ike_version"])
    d.Set("encryption_algorithm", data["encryption_algorithm"])
    d.Set("pfs", data["pfs"])
    d.Set("lifetime", data["lifetime"])
    d.Set("user", data["user"])

    return nil
}
