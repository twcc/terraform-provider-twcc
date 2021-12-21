package twcc

import (
    "fmt"
    "log"
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceVolumeSnapshot() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceVolumeSnapshotRead,

        Schema: map[string]*schema.Schema{
            "create_time": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "desc": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "name": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "project": {
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

            "volume": {
                Type:		schema.TypeString,
                Computed:	true,
            },
        },
    }
}

// dataSourceVolumeSnapshotRead performs the volume snapshot lookup.
func dataSourceVolumeSnapshotRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    platform := d.Get("platform").(string)
    name := d.Get("name").(string)
    projectID := d.Get("project").(string)

    resourcePath := fmt.Sprintf("api/v3/%s/snapshots/?project=%s", platform, projectID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to list snapshots: %v", err)
    }

    var data []map[string]interface{}
    if err = json.Unmarshal([]byte(response), &data); err != nil {
        return err
    }

    for _, snapshot := range data {
        if snapshot["name"] == name {
            return dataSourceVolumeSnapshotAttributes(d, snapshot)
        }
    }

    return fmt.Errorf("Unable to retrieve volume snapshot %s: %v", name, err)
}

// dataSourceVolumeSnapshotAttributes populates the fields of a volume snapshot data source.
func dataSourceVolumeSnapshotAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    snapshot_id := int(data["id"].(float64))
    log.Printf("[DEBUG] Retrieved twcc_volume_snapshot: %d", snapshot_id)

    d.SetId(fmt.Sprintf("%d", snapshot_id))
    d.Set("create_time", data["create_time"])
    d.Set("desc", data["desc"])
    d.Set("status", data["status"])
    d.Set("user", data["user"])
    volumeInfo := data["volume"].(map[string]interface{})
    volumeID := int(volumeInfo["id"].(float64))
    d.Set("volume", fmt.Sprintf("%d", volumeID))

    return nil
}
