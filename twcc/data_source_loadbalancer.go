package twcc

import (
    "fmt"
    "log"
    "encoding/json"
    "strings"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceLoadBalancer() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceLoadBalancerRead,

        Schema: map[string]*schema.Schema{
            "create_time": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "desc": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "listeners": {
                Type:		schema.TypeList,
                Computed:	true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "default_tls_container_ref": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },

                        "name": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },

                        "pool": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },

                        "protocol": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },

                        "protocol_port": {
                            Type:	schema.TypeInt,
                            Computed:	true,
                        },

                        "sni_container_refs": {
                            Type:	schema.TypeList,
                            Computed:	true,
                            Elem:	&schema.Schema{
                                Type:	schema.TypeString,
                            },
                        },

                        "status": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },
                    },
                },
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

            "pools": {
                Type:		schema.TypeList,
                Computed:	true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "id": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },

                        "method": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },

                        "name": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },

                        "protocol": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },

                        "status": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },
                    },
                },
            },

            "private_net": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "project": {
                Type:		schema.TypeString,
                Optional:	true,
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

// dataSourceLoadBalancerRead performs the loadbalancer lookup.
func dataSourceLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    name := d.Get("name").(string)
    platform := d.Get("platform").(string)
    params := []string{fmt.Sprintf("name=%s", name)}
    if project := d.Get("project"); project != "" {
        params = append(params, fmt.Sprintf("project=%s", project))
    }
    if private_net := d.Get("private_net"); private_net != "" {
        params = append(params, fmt.Sprintf("private_net=%s", private_net))
    }
    if len(params) < 2 {
        return fmt.Errorf("Either project or private_net should be defined")
    }

    resourcePath := fmt.Sprintf("api/v3/%s/loadbalancers/?%s", platform,
                                strings.Join(params, "&"))
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to list loadbalancers: %v", err)
    }

    var data []map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    for _, loadbalancer := range data {
        if loadbalancer["name"] == name {
            return dataSourceLoadBalancerAttributes(d, loadbalancer)
        }
    }

    return fmt.Errorf("Unable to retrieve loadbalancer %s: %v", name, err)
}

// dataSourceLoadBalancerAttributes populates the fields of a loadbalancer data source.
func dataSourceLoadBalancerAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    loadbalancer_id := int(data["id"].(float64))
    log.Printf("[DEBUG] Retrieved twcc_loadbalancer: %d", loadbalancer_id)

    d.SetId(fmt.Sprintf("%d", loadbalancer_id))
    d.Set("status", data["status"])
    d.Set("user", data["user"])
    d.Set("name", data["name"])
    d.Set("desc", data["desc"])
    private_net_info := data["private_net"].(map[string]interface{})
    private_net := fmt.Sprintf("%d", int(private_net_info["id"].(float64)))
    d.Set("private_net", private_net)
    listeners := data["listeners"].([]interface{})
    listenerArray := make([]interface{}, len(listeners))
    for i, listener := range listeners {
        data := listener.(map[string]interface{})
        info := make(map[string]interface{})
        if dtcr, ok := data["default_tls_container_ref"].(float64); ok {
            info["default_tls_container_ref"] = fmt.Sprintf("%d", int(dtcr))
        } else {
            info["default_tls_container_ref"] = ""
        }
        info["name"] = data["name"].(string)
        info["pool"] = fmt.Sprintf("%d", int(data["pool"].(float64)))
        info["protocol"] = data["protocol"].(string)
        info["protocol_port"] = int(data["protocol_port"].(float64))
        scrs := data["sni_container_refs"].([]interface{})
        scrArray := make([]string, len(scrs))
        for j, scr := range scrs {
            scrArray[j] = fmt.Sprintf("%d", int(scr.(float64)))
        }
        info["sni_container_refs"] = scrArray
        info["status"] = data["status"].(string)
        listenerArray[i] = info
    }
    d.Set("listeners", listenerArray)
    pools := data["pools"].([]interface{})
    poolArray := make([]interface{}, len(pools))
    for i, pool := range pools {
        data := pool.(map[string]interface{})
        info := make(map[string]interface{})
        info["id"] = fmt.Sprintf("%d", int(data["id"].(float64)))
        info["method"] = data["method"].(string)
        info["name"] = data["name"].(string)
        info["protocol"] = data["protocol"].(string)
        info["status"] = data["status"].(string)
        poolArray[i] = info
    }
    d.Set("pools", poolArray)

    return nil
}
