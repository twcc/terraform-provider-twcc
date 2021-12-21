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

type VolumeSnapshotCreateBody struct {
    Desc		string	`json:"desc,omitempty"`
    Name		string	`json:"name"`
    Volume		string	`json:"volume"`
}

func resourceVolumeSnapshot() *schema.Resource {
    return &schema.Resource{
        Create: resourceVolumeSnapshotCreate,
        Read:   resourceVolumeSnapshotRead,
        Delete: resourceVolumeSnapshotDelete,

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

            "restore_volume": {     
                Type:		schema.TypeString,
                Computed:	true,
            },

            "snapshot_uuid": {
                Type:		schema.TypeString,
                Computed:	true,
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

            "volume": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },
        },
    }
}

func resourceVolumeSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    name := d.Get("name").(string)
    platform := d.Get("platform").(string)
    desc := d.Get("desc").(string)
    volume := d.Get("volume").(string)
    body := VolumeSnapshotCreateBody {
        Desc:	desc,
        Name:	name,
        Volume:	volume,
    }

    resourcePath := fmt.Sprintf("api/v3/%s/snapshots/", platform)
    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    response, err := config.doNormalRequest(platform, resourcePath, "POST", buf)

    if err != nil {
        return fmt.Errorf("Error creating twcc_snapshot %s on %s: %v", name, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    snapshotID := int(data["id"].(float64))
    d.SetId(fmt.Sprintf("%d", snapshotID))

    newPath := fmt.Sprintf("%s/%d/", resourcePath, snapshotID)
    stateConf := &resource.StateChangeConf{
        Pending:    []string{"CREATING",},
        Target:     []string{"AVAILABLE", "ERROR"},
        Refresh:    snapshotStateRefreshFunc(config, platform, newPath),
        Timeout:    d.Timeout(schema.TimeoutCreate),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf(
            "Error waiting for twcc_volume_snapshot %d to become AVAILABLE: %v", snapshotID, err)
    }

    d.Set("name", name)
    d.Set("platform", platform)
    d.Set("volume", volume)
    d.Set("desc", desc)
    return resourceVolumeSnapshotRead(d, meta)
}

func resourceVolumeSnapshotRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    snapshotID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/snapshots/%s/", platform, snapshotID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to retrieve volume snapshot %s on %s: %v", snapshotID, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve volume snapshot json data: %s", err)
    }

    log.Printf("[DEBUG] Retrieved twcc_volume_snapshot %s", d.Id())
    if hostList, ok := data["attached_host"].([]interface{}); ok {
        hostInfo := flattenVolumeHostInfo(hostList)
        d.Set("attached_host", hostInfo)
    } else {
        var emptyList []interface{}
         d.Set("attached_host", emptyList)
    }

    d.Set("create_time", data["create_time"])
    d.Set("restore_volume", data["restore_volume"])
    d.Set("snapshot_uuid", data["snapshot_uuid"])
    d.Set("status", data["status"])
    d.Set("status_reason", data["status_reason"])
    d.Set("user", data["user"].(map[string]interface{}))
    return nil
}

func resourceVolumeSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    platform := d.Get("platform").(string)
    snapshotID := d.Id()
    resourcePath := fmt.Sprintf("api/v3/%s/snapshots/%s/", platform, snapshotID)
    _, err := config.doNormalRequest(platform, resourcePath, "DELETE", nil)

    if err != nil {
        return fmt.Errorf("Unable to delete volume snapshot %s: on %s %v", snapshotID, platform, err)
    }

    stateConf := &resource.StateChangeConf{
        Pending:    []string{"DELETING",},
        Target:     []string{"DELETED", "ERROR"},
        Refresh:    snapshotStateRefreshForDeletedFunc(config, platform, resourcePath),
        Timeout:    d.Timeout(schema.TimeoutDelete),
        Delay:      10 * time.Second,
        MinTimeout: 3 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf( 
            "Error waiting for twcc_volume_snapshot %s to become DELETED: %v", snapshotID, err)
    }

    d.SetId("")

    return nil
}
