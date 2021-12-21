package twcc

import (
    "fmt"
    "log"
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceVolume() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceVolumeRead,

        Schema: map[string]*schema.Schema{
            "attached_host": {
                Type:		schema.TypeMap,
                Computed:	true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },

            "create_time": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "is_attached": {
                Type:		schema.TypeBool,
                Computed:	true,
            },

            "is_bootable": {
                Type:		schema.TypeBool,
                Computed:	true,
            },

            "mountpoint": {
                Type:		schema.TypeList,
                Computed:	true,
                Elem: &schema.Schema{
                    Type:	schema.TypeString,
                },
            },

            "platform": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "name": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "project": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "vcs": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "size": {
                Type:		schema.TypeInt,
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

            "volume_type": {
                Type:		schema.TypeString,
                Computed:	true,
            },
        },
    }
}

// dataSourceVolumeRead performs the volume lookup.
func dataSourceVolumeRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    platform := d.Get("platform").(string)
    serversID := []int{}
    var name, siteID, projectID string
    if site := d.Get("vcs"); site != "" {
        siteID = site.(string)
    }
    if volName := d.Get("name"); volName != "" {
        name = volName.(string)
    }
    if project := d.Get("project"); project != "" {
        projectID = project.(string)
    }

    if siteID != "" {
        resourcePath := fmt.Sprintf("api/v3/%s/sites/%s/", platform, siteID)
        response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)
        if err != nil {
            return fmt.Errorf("Unable to get VCS %s: %v", siteID, err)
        }
        var data map[string]interface{}
        if err = json.Unmarshal([]byte(response), &data); err != nil {
            return err
        }
        projectID = fmt.Sprintf("%d", int(data["project"].(float64)))
        for _, server := range data["servers"].([]interface{}) {
            serverInfo := server.(map[string]interface{})
            serversID = append(serversID, int(serverInfo["id"].(float64)))
        }
    } else if name == "" || projectID == "" {
        return fmt.Errorf("name and project are required when vcs is not defined")
    }

    resourcePath := fmt.Sprintf("api/v3/%s/volumes/?project=%s", platform, projectID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to list volumes: %v", err)
    }

    var data []map[string]interface{}
    if err = json.Unmarshal([]byte(response), &data); err != nil {
        return err
    }

    found := false
    for _, volume := range data {
        if len(serversID) > 0 {
            if host, ok := volume["attached_host"].(map[string]interface{}); ok {
                attached_hostID := int(host["id"].(float64))
                for _, serverID := range serversID {
                    if attached_hostID == serverID {
                        found = true
                    }
                }
            }
        }
        if name != "" {
            if volume["name"] == name {
                found = true
            } else if found {
                found = false
            }
        }
        if found {
            return dataSourceVolumeAttributes(d, volume)
        }
    }
    if siteID != "" {
        return fmt.Errorf("Unable to retrieve volume by VCS %s: %v", siteID, err)
    } else {
        return fmt.Errorf("Unable to retrieve volume %s: %v", name, err)
    }
}

// dataSourceVolumeAttributes populates the fields of a volume data source.
func dataSourceVolumeAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    volume_id := int(data["id"].(float64))
    log.Printf("[DEBUG] Retrieved twcc_volume: %d", volume_id)

    d.SetId(fmt.Sprintf("%d", volume_id))
    if host, ok := data["attached_host"].(interface{}); ok {
        hostInfo := flattenVolumeHostInfo(host)
        d.Set("attached_host", hostInfo)
    } else {
        var emptyHost interface{}
        d.Set("attached_host", emptyHost)
    }
    d.Set("create_time", data["create_time"])
    d.Set("is_attached", data["is_attached"])
    d.Set("is_bootable", data["is_bootable"])
    d.Set("mountpoint", data["mountpoint"])
    d.Set("name", data["name"])
    projectInfo := data["project"].(map[string]interface{})
    projectID := int(projectInfo["id"].(float64))
    d.Set("project", fmt.Sprintf("%d", projectID))
    d.Set("size", data["size"])
    d.Set("status", data["status"])
    d.Set("volume_type", data["volume_type"])
    d.Set("user", data["user"])

    return nil
}
