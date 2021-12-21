package twcc

import (
    "fmt"
    "log"
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceContainer() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceContainerRead,

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

            "pod": {
                Type:           schema.TypeList,
                Computed:       true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "container": {
                            Type:	schema.TypeList,
                            Computed:	true,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "image": {
                                        Type:		schema.TypeString,
                                        Computed:	true,
                                    },
                                    "name": {
                                        Type:		schema.TypeString,
                                        Computed:	true,
                                    },
                                    "ports": {
                                        Type:		schema.TypeList,
                                        Computed:	true,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "name": {
                                                    Type:	schema.TypeString,
                                                    Computed:	true,
                                                },

                                                "port": {
                                                    Type:	schema.TypeInt,
                                                    Computed:	true,
                                                },

                                                "protocol": {
                                                    Type:	schema.TypeString,
                                                    Computed:	true,
                                                },
                                            },
                                        },
                                    },
                                    "volumes": {
                                        Type:		schema.TypeList,
                                        Computed:	true,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "mount_path": {
                                                    Type:	schema.TypeString,
                                                    Computed:	true,
                                                },

                                                "path": {
                                                    Type:	schema.TypeString,
                                                    Computed:	true,
                                                },

                                                "read_only": {
                                                    Type:	schema.TypeBool,
                                                    Computed:	true,
                                                },

                                                "type": {
                                                    Type:	schema.TypeString,
                                                    Computed:	true,
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },

                        "flavor": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },

                        "message": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },

                        "name": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },

                        "reason": {
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

            "service": {
                Type:		schema.TypeList,
                Computed:	true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "name": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },

                        "net_type": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },

                        "ports": {
                            Type:	schema.TypeList,
                            Computed:	true,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "port": {
                                        Type:		schema.TypeInt,
                                        Computed:	true,
                                    },

                                    "protocol": {
                                        Type:		schema.TypeString,
                                        Computed:	true,
                                    },

                                    "target_port": {
                                        Type:		schema.TypeInt,
                                        Computed:	true,
                                    },
                                },
                            },
                        },

                        "public_ip": {
                            Type:	schema.TypeList,
                            Computed:	true,
                            Elem: &schema.Schema{
                                Type:	schema.TypeString,
                            },
                        },
                    },
                },
            },
        },
    }
}

// dataSourceContainerRead performs the container lookup.
func dataSourceContainerRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    name := d.Get("name").(string)
    platform := d.Get("platform").(string)
    projectID := d.Get("project").(string)

    resourcePath := fmt.Sprintf("api/v3/%s/sites/?project=%s&name=%s", platform, projectID, name)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to list containers: %v", err)
    }

    var data []map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    var site_id int
    var siteInfo, containerInfo map[string]interface{}
    for _, site := range data {
        if site["name"] == name {
            if site_id != 0 {
                return fmt.Errorf("There are duplicated containers with name '%s'", name)
            }
            site_id = int(site["id"].(float64))
            siteInfo = site
        }
    }
    if site_id != 0 {
        resourcePath := fmt.Sprintf("api/v3/%s/sites/%d/container/", platform, site_id)
        response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)
        if err != nil {
            return fmt.Errorf("Unable to retrieve container: %v", err)
        }

        if err = json.Unmarshal([]byte(response), &containerInfo); err != nil {
            return err
        }
        return dataSourceContainerAttributes(d, siteInfo, containerInfo)
    }

    return fmt.Errorf("Unable to retrieve container %s: %v", name, err)
}

// dataSourceContainerAttributes populates the fields of a container data source.
func dataSourceContainerAttributes(d *schema.ResourceData, siteInfo map[string]interface{}, containerInfo map[string]interface{}) error {
    site_id := int(siteInfo["id"].(float64))
    log.Printf("[DEBUG] Retrieved twcc_container: %d", site_id)

    d.SetId(fmt.Sprintf("%d", site_id))
    d.Set("create_time", siteInfo["create_time"])
    d.Set("public_ip", siteInfo["public_ip"])
    d.Set("solution", fmt.Sprintf("%d", int(siteInfo["solution"].(float64))))
    d.Set("pod", flattenSitePodInfo(containerInfo["Pod"].([]interface{})))
    d.Set("service", flattenSiteServiceInfo(containerInfo["Service"].([]interface{})))
    d.Set("status", siteInfo["status"])
    d.Set("user", siteInfo["user"])

    return nil
}
