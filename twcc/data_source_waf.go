package twcc

import (
    "fmt"
    "log"
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceWAF() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceWAFRead,

        Schema: map[string]*schema.Schema{
            "create_time": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "name": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:   true,
            },

            "platform": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:   true,
            },

            "project": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "public_ip": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "servers": {
                Type:		schema.TypeList,
                Computed:	true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "flavor_id": {
                            Type:		schema.TypeString,
                            Computed:	true,
                        },
                        "hostname": {
                            Type:		schema.TypeString,
                            Computed:	true,
                        },

                        "id": {
                            Type:		schema.TypeString,
                            Computed:	true,
                        },

                        "status": {
                            Type:		schema.TypeString,
                            Computed:	true,
                        },
                    },
                },
            },

            "solution": {
                Type:		schema.TypeString,
                Computed:	true,
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

// dataSourceWAFRead performs the WAF lookup.
func dataSourceWAFRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    name := d.Get("name").(string)
    platform := d.Get("platform").(string)
    projectID := d.Get("project").(string)

    resourcePath := fmt.Sprintf("api/v3/%s/sites/?project=%s&name=%s", platform, projectID, name)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to list WAF: %v", err)
    }

    var data []map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    found := false
    for _, waf := range data {
        if waf["name"] == name {
            if found {
                return fmt.Errorf("There are duplicated WAF with name '%s'", name)
            }
            err = dataSourceWAFAttributes(d, waf)
            found = true
        }
    }
    if found {
        return err
    }

    return fmt.Errorf("Unable to retrieve WAF %s: %v", name, err)
}

// dataSourceWAFAttributes populates the fields of a WAF data source.
func dataSourceWAFAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    waf_id := int(data["id"].(float64))
    log.Printf("[DEBUG] Retrieved twcc_waf: %d", waf_id)

    d.SetId(fmt.Sprintf("%d", waf_id))
    d.Set("create_time", data["create_time"])
    d.Set("public_ip", data["public_ip"])
    d.Set("solution", fmt.Sprintf("%d", int(data["solution"].(float64))))
    serversInfo := flattenSiteServersInfo(data["servers"].([]interface{}))
    d.Set("servers", serversInfo)
    d.Set("status", data["status"])
    d.Set("user", data["user"])

    return nil
}
