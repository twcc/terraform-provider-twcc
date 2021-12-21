package twcc

import (
    "fmt"
    "log"
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceFirewallRule() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceFirewallRuleRead,

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

            "action": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "destination_ip_address": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "destination_port": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "source_ip_address": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "source_port": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "ip_version": {
                Type:		schema.TypeInt,
                Computed:	true,
            },

            "protocol": {
                Type:		schema.TypeString,
                Computed:	true,
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

// dataSourceFirewallRuleRead performs the firewall rule lookup.
func dataSourceFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    name := d.Get("name").(string)
    projectID := d.Get("project").(string)
    platform := d.Get("platform").(string)

    resourcePath := fmt.Sprintf("api/v3/%s/firewall_rules/?project=%s", platform, projectID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to list firewall rules: %v", err)
    }

    var data []map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    for _, firewall_rule := range data {
        if firewall_rule["name"] == name {
            return dataSourceFirewallRuleAttributes(d, firewall_rule)
        }
    }

    return fmt.Errorf("Unable to retrieve firewall rule %s: %v", name, err)
}

// dataSourceFirewallRuleAttributes populates the fields of a firewall rule data source.
func dataSourceFirewallRuleAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    firewall_rule_id := int(data["id"].(float64))
    log.Printf("[DEBUG] Retrieved twcc_firewall_rule: %d", firewall_rule_id)

    d.SetId(fmt.Sprintf("%d", firewall_rule_id))
    d.Set("project", fmt.Sprintf("%v", data["project"]))
    d.Set("user", data["user"])
    d.Set("name", data["name"])
    d.Set("platform", data["platform"])
    d.Set("protocol", data["protocol"])
    d.Set("ip_version", data["ip_version"])
    d.Set("action", data["action"])
    d.Set("destination_ip_address", data["destination_ip_address"])
    d.Set("destination_port", data["destination_port"])
    d.Set("source_ip_address", data["source_ip_address"])
    d.Set("source_port", data["source_port"])

    return nil
}
