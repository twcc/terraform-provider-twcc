package twcc

import (
    "fmt"
    "log"
    "encoding/json"
    "strings"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceVPN() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceVPNRead,

        Schema: map[string]*schema.Schema{
            "ike_policy": {
                Type:		schema.TypeString,
                Optional:	true,
                Computed:	true,
                ForceNew:	true,
            },

            "ipsec_policy": {
                Type:		schema.TypeString,
                Optional:	true,
                Computed:	true,
                ForceNew:	true,
            },

            "private_network": {
                Type:		schema.TypeString,
                Optional:	true,
                Computed:	true,
                ForceNew:	true,
            },

            "project": {
                Type:		schema.TypeString,
                Optional:	true,
                Computed:	true,
                ForceNew:	true,
            },

            "local_address": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "local_cidr": {
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

            "vpn_connection": {
                Type:		schema.TypeList,
                Computed:	true,
                Elem:		&schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "dpd_action": {
                            Type:		schema.TypeString,
                            Computed:	true,
                        },

                        "dpd_interval": {
                            Type:		schema.TypeInt,
                            Computed:	true,
                        },

                        "dpd_timeout": {
                            Type:		schema.TypeInt,
                            Computed:	true,
                        },

                        "initiator": {
                            Type:		schema.TypeString,
                            Computed:	true,
                        },

                        "mtu": {
                            Type:		schema.TypeInt,
                            Computed:	true,
                        },

                        "peer_address": {
                            Type:		schema.TypeString,
                            Computed:	true,
                        },

                        "peer_cidrs": {
                            Type:		schema.TypeList,
                            Computed:	true,
                            Elem: &schema.Schema{
                                Type:	schema.TypeString,
                            },
                        },

                        "peer_id": {
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
        },
    }
}

// dataSourceVPNRead performs the VPN lookup.
func dataSourceVPNRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    name := d.Get("name").(string)
    platform := d.Get("platform").(string)
    var params []string
    if project := d.Get("project"); project != "" {
        params = append(params, fmt.Sprintf("project=%s", project))
    }
    if ike_policy := d.Get("ike_policy"); ike_policy != "" {
        params = append(params, fmt.Sprintf("ike_policy=%s", ike_policy))
    }
    if ipsec_policy := d.Get("ipsec_policy"); ipsec_policy != "" {
        params = append(params, fmt.Sprintf("ipsec_policy=%s", ipsec_policy))
    }
    if private_network := d.Get("private_network"); private_network != "" {
        params = append(params, fmt.Sprintf("private_network=%s", private_network))
    }

    resourcePath := fmt.Sprintf("api/v3/%s/vpn_services/?%s",
                                platform, strings.Join(params, "&"))
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to list vpn: %v", err)
    }

    var data []map[string]interface{}
    if err = json.Unmarshal([]byte(response), &data); err != nil {
        return err
    }

    var vpn_id int
    for _, vpn := range data {
        if vpn["name"] == name {
            if vpn_id != 0 {
                return fmt.Errorf("There are duplicated vpn with name '%s'", name)
            }
            vpn_id = int(vpn["id"].(float64))
        }
    }
    if vpn_id != 0 {
        resourcePath := fmt.Sprintf("api/v3/%s/vpn_services/%d/", platform, vpn_id)
        response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)
        if err != nil {
            return fmt.Errorf("Unable to retrieve vpn: %v", err)
        }

        var vpn map[string]interface{}
        if err = json.Unmarshal([]byte(response), &vpn); err != nil {
            return err
        }
        return dataSourceVPNAttributes(d, vpn)
    }

    return fmt.Errorf("Unable to retrieve vpn %s: %v", name, err)
}

// retrieve id from a map
func retrieveID(data interface{}) string {
    dataMap := data.(map[string]interface{})
    ID := int(dataMap["id"].(float64))
    return fmt.Sprintf("%d", ID)
}

// dataSourceVPNAttributes populates the fields of a VPN data source.
func dataSourceVPNAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    vpn_id := int(data["id"].(float64))
    log.Printf("[DEBUG] Retrieved twcc_vpn: %d", vpn_id)

    d.SetId(fmt.Sprintf("%d", vpn_id))
    d.Set("user", data["user"])
    d.Set("local_address", data["local_address"])
    d.Set("local_cidr", data["local_cidr"])
    d.Set("status", data["status"])
    d.Set("ike_policy", retrieveID(data["ike_policy"]))
    d.Set("ipsec_policy", retrieveID(data["ipsec_policy"]))
    d.Set("private_network", retrieveID(data["private_network"]))
    d.Set("vpn_connection", flattenVPNConnectionInfo(data["connection"]))

    return nil
}
