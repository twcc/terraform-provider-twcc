package twcc

import (
    "fmt"
    "log"
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSecurityGroup() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceSecurityGroupRead,

        Schema: map[string]*schema.Schema{
            "platform": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "vcs": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "name": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "security_group_rules": {
                Type:		schema.TypeList,
                Computed:	true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "id": {
                            Type:		schema.TypeString,
                            Computed:	true,
                        },
                        "direction": {
                            Type:		schema.TypeString,
                            Computed:	true,
                        },
                        "ethertype": {
                            Type:		schema.TypeString,
                            Computed:	true,
                        },
                        "protocol": {
                            Type:		schema.TypeString,
                            Computed:	true,
                        },
                        "remote_ip_prefix": {
                            Type:		schema.TypeString,
                            Computed:	true,
                        },
                        "port_range_min": {
                            Type:		schema.TypeInt,
                            Computed:	true,
                        },
                        "port_range_max": {
                            Type:		schema.TypeInt,
                            Computed:	true,
                        },
                    },
                },
            },
        },
    }
}

// dataSourceSecurityGroupRead performs the security group lookup.
func dataSourceSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    platform := d.Get("platform").(string)
    siteID := d.Get("vcs").(string)

    resourcePath := fmt.Sprintf("api/v3/%s/sites/%s/", platform, siteID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)
    if err != nil {
        return fmt.Errorf("Unable to get VCS %s: %v", siteID, err)
    }
    var site map[string]interface{}
    if err = json.Unmarshal([]byte(response), &site); err != nil {
        return err
    }
    projectID := int(site["project"].(float64))
    var serverID int
    for _, server := range site["servers"].([]interface{}) {
        serverInfo := server.(map[string]interface{})
        serverID = int(serverInfo["id"].(float64))
        break
    }

    resourcePath = fmt.Sprintf("api/v3/%s/security_groups/?project=%d&server=%d", platform, projectID, serverID)
    response, err = config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to list security_groups: %v", err)
    }

    var security_groups []map[string]interface{}
    if err = json.Unmarshal([]byte(response), &security_groups); err != nil {
        return err
    }

    for _, security_group := range security_groups {
        return dataSourceSecurityGroupAttributes(d, security_group)
    }

    return fmt.Errorf("Unable to retrieve security group by VCS %s: %v", siteID, err)
}

// dataSourceSecurityGroupAttributes populates the fields of a security group data source.
func dataSourceSecurityGroupAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    security_group_id := data["id"].(string)
    log.Printf("[DEBUG] Retrieved twcc_security_group: %s", security_group_id)

    d.SetId(security_group_id)
    d.Set("name", data["name"])
    security_group_rules := flattenSecurityGroupRulesInfo(data["security_group_rules"].([]interface{}))
    d.Set("security_group_rules", security_group_rules)

    return nil
}
