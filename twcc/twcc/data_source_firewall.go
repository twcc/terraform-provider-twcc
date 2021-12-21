package twcc

import (
    "fmt"
    "log"
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceFirewall() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceFirewallRead,

        Schema: map[string]*schema.Schema{
            "name": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:   true,
            },

            "platform": {
                Type:       schema.TypeString,
                Required:   true,
                ForceNew:   true,
            },

            "project": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "status": {
                Type:		schema.TypeString,
                Computed:   true,
            },

            "desc": {
                Type:		schema.TypeString,
                Computed:   true,
            },

            "user": {
                Type:       schema.TypeMap,
                Computed:   true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
        },
    }
}

// dataSourceFirewallRead performs the firewall lookup.
func dataSourceFirewallRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    name := d.Get("name").(string)
    projectID := d.Get("project").(string)
    platform := d.Get("platform").(string)

    resourcePath := fmt.Sprintf("api/v3/%s/firewalls/?project=%s", platform, projectID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to list firewalls: %v", err)
    }

    var data []map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    for _, firewall := range data {
        if firewall["name"] == name {
            return dataSourceFirewallAttributes(d, firewall)
        }
    }

    return fmt.Errorf("Unable to retrieve firewall %s: %v", name, err)
}

// dataSourceFirewallAttributes populates the fields of a firewall data source.
func dataSourceFirewallAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    firewall_id := int(data["id"].(float64))
    log.Printf("[DEBUG] Retrieved twcc_firewall: %d", firewall_id)

    d.SetId(fmt.Sprintf("%d", firewall_id))
    d.Set("project", fmt.Sprintf("%v", data["project"]))
    d.Set("status", data["status"])
    d.Set("user", data["user"])
    d.Set("name", data["name"])
    d.Set("desc", data["desc"])
    d.Set("platform", data["platform"])

    return nil
}
