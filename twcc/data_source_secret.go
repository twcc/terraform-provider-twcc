package twcc

import (
    "fmt"
    "log"
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSecret() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceSecretRead,

        Schema: map[string]*schema.Schema{
            "create_time": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "desc": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "expire_time": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "name": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "project": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "platform": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "status": {
                Type:		schema.TypeString,
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

// dataSourceSecretRead performs the secret lookup.
func dataSourceSecretRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    platform := d.Get("platform").(string)
    name := d.Get("name").(string)
    projectID := d.Get("project").(string)

    resourcePath := fmt.Sprintf("api/v3/%s/secrets/?project=%s", platform, projectID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to list secrets: %v", err)
    }

    var data []map[string]interface{}
    if err = json.Unmarshal([]byte(response), &data); err != nil {
        return err
    }

    for _, secret := range data {
        if secret["name"] == name {
            return dataSourceSecretAttributes(d, secret)
        }
    }

    return fmt.Errorf("Unable to retrieve secret %s: %v", name, err)
}

// dataSourceSecretAttributes populates the fields of a secret data source.
func dataSourceSecretAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    secret_id := int(data["id"].(float64))
    log.Printf("[DEBUG] Retrieved twcc_secret: %d", secret_id)

    d.SetId(fmt.Sprintf("%d", secret_id))
    d.Set("create_time", data["create_time"])
    d.Set("desc", data["desc"])
    d.Set("expire_time", data["expire_time"])
    d.Set("status", data["status"])
    d.Set("user", data["user"])

    return nil
}
