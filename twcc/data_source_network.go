package twcc

import (
    "fmt"
    "log"
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNetwork() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceNetworkRead,

        Schema: map[string]*schema.Schema{
            "cidr": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "create_time": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "dns_domain": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "ext_net": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "gateway": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "ip_version": {
                Type:		schema.TypeInt,
                Computed:	true,
            },

            "name": {
                Type:		schema.TypeString,
                Required:	true,
            },

            "nameservers": {
                Type:		schema.TypeList,
                Computed:	true,
                Elem: &schema.Schema{
                    Type:	schema.TypeString,
                },
            },

            "platform": {
                Type:		schema.TypeString,
                Required:	true,
            },

            "project": {
                Type:		schema.TypeString,
                Required:	true,
            },

            "status": {
                Type:		schema.TypeString,
                Computed:       true,
            },

            "user": {
                Type:           schema.TypeMap,
                Computed:       true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },

            "with_router": {
                Type:           schema.TypeBool,
                Optional:       true,
                ForceNew:       true,
            },
        },
    }
}

// dataSourceNetworkRead performs the network lookup.
func dataSourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    name := d.Get("name").(string)
    projectID := d.Get("project").(string)
    platform := d.Get("platform").(string)

    resourcePath := fmt.Sprintf("api/v3/%s/networks/?project=%s", platform, projectID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to list networks: %v", err)
    }

    var data []map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    for i := range data {
        if data[i]["name"] == name {
            return dataSourceNetworkAttributes(d, data[i])
        }
    }

    return fmt.Errorf("Unable to retrieve network %s: %v", name, err)
}

// dataSourceNetworkAttributes populates the fields of a network data source.
func dataSourceNetworkAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    network_id := int(data["id"].(float64))
    log.Printf("[DEBUG] Retrieved twcc_network: %d", network_id)

    d.SetId(fmt.Sprintf("%d", network_id))
    d.Set("cidr", data["cidr"])
    d.Set("create_time", data["create_time"])
    d.Set("dns_domain", data["dns_domain"])
    d.Set("ext_net", data["ext_net"])
    d.Set("gateway", data["gateway"])
    d.Set("ip_version", data["ip_version"])
    d.Set("nameservers", data["nameservers"])
    d.Set("project", fmt.Sprintf("%v", data["project"]))
    d.Set("status", data["status"])
    d.Set("user", data["user"])
    d.Set("with_router", data["with_router"])
    d.Set("name", data["name"])
    d.Set("platform", data["platform"])

    return nil
}
