package twcc

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ContainerCreateBody struct {
    Desc	string	`json:"desc,omitempty"`
    Name	string	`json:"name"`
    Project	string	`json:"project"`
    Solution	string	`json:"solution"`
}

func resourceContainer() *schema.Resource {
    return &schema.Resource{
        Create: resourceContainerCreate,
        Read:   resourceContainerRead,
        Delete: resourceContainerDelete,

        Timeouts: &schema.ResourceTimeout{
            Create: schema.DefaultTimeout(10 * time.Minute),
            Delete: schema.DefaultTimeout(10 * time.Minute),
        },

        Schema: map[string]*schema.Schema{
            "create_time": {
                Type:		schema.TypeString,
                Computed:	true,
                ForceNew:	true,
            },

            "desc": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "extra_property": {
                Type:		schema.TypeMap,
                Optional:	true,
                ForceNew:	true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },

            "name": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:       true,
            },

            "platform": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:       true,
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

            "project": {     
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "public_ip": {
                Type:		schema.TypeString,
                Computed:	true,
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

            "solution": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "status": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "status_reason": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "user": {       
                Type:		schema.TypeMap,
                Computed:	true,
                ForceNew:	true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
        },
    }
}

func resourceContainerCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    desc := d.Get("desc").(string)
    extra_property := d.Get("extra_property").(map[string]interface{})
    headers := make(map[string]string)
    for key, value := range extra_property {
        header := fmt.Sprintf("x-extra-property-%s", key)
        headers[header] = fmt.Sprintf("%v", value)
    }
    name := d.Get("name").(string)
    platform := d.Get("platform").(string)
    project := d.Get("project").(string)
    solution := d.Get("solution").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/sites/", platform)

    body := ContainerCreateBody {
        Desc:		desc,
        Name:		name,
        Project:	project,
        Solution:	solution,
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    response, err := config.doCreateSiteRequest(platform, resourcePath, "POST", buf, headers)

    if err != nil {
        return fmt.Errorf("Error creating twcc_container %s on %s: %v", name, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    siteID := int(data["id"].(float64))
    d.SetId(fmt.Sprintf("%d", siteID))

    newPath := fmt.Sprintf("%s%d/", resourcePath, siteID)
    stateConf := &resource.StateChangeConf{
        Pending:    []string{"Error", "Initializing", "Queueing"},
        Target:     []string{"Ready"},
        Refresh:    siteStateRefreshFunc(config, platform, newPath),
        Timeout:    d.Timeout(schema.TimeoutCreate),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf(
            "Error waiting for twcc_container %d to become Ready: %v", siteID, err)
    }

    d.Set("name", name)
    d.Set("platform", platform)
    d.Set("project", project)
    d.Set("solution", solution)
    return resourceContainerRead(d, meta)
}

func resourceContainerRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    siteID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/sites/%s/", platform, siteID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to retrieve container %s on %s: %v", siteID, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve container json data: %v", err)
    }

    log.Printf("[DEBUG] Retrieved twcc_container %s", d.Id())
    d.Set("create_time", data["create_time"])
    d.Set("desc", data["desc"])
    d.Set("public_ip", data["public_ip"])
    d.Set("status", data["status"])
    d.Set("status_reason", data["status_reason"])
    d.Set("user", data["user"])

    containerPath := fmt.Sprintf("%s/container/", resourcePath)
    response, err = config.doNormalRequest(platform, containerPath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to retrieve container %s detail on %s: %v", siteID, platform, err)
    }

    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve container detail json data: %v", err)
    }

    log.Printf("[DEBUG] Retrieved twcc_container detail %s", d.Id())
    podInfo := flattenSitePodInfo(data["Pod"].([]interface{}))
    d.Set("pod", podInfo)
    serviceInfo := flattenSiteServiceInfo(data["Service"].([]interface{}))
    d.Set("service", serviceInfo)
    return nil
}

func resourceContainerDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    siteID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/sites/%s/", platform, siteID)
    _, err := config.doNormalRequest(platform, resourcePath, "DELETE", nil)

    if err != nil {
        return fmt.Errorf("Unable to delete container %s: on %s %v", siteID, platform, err)
    }

    stateConf := &resource.StateChangeConf{
        Pending:    []string{"Deleting"},
        Target:     []string{"Deleted"},
        Refresh:    siteStateRefreshForDeletedFunc(config, platform, resourcePath),
        Timeout:    d.Timeout(schema.TimeoutDelete),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf( 
            "Error waiting for twcc_container %s to become Deleted: %v", siteID, err)
    }

    d.SetId("")

    return nil
}
