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

type VolumeCreateBody struct {
    Desc		string	`json:"desc,omitempty"`
    Name		string	`json:"name"`
    Project		string	`json:"project,omitempty"`
    Size		int	`json:"size,omitempty"`
    SrcSnapshot		string	`json:"src_snapshot,omitempty"`
    VolumeType		string	`json:"volume_type,omitempty"`
}

type VolumeUpdateBody struct {
    Status	string	`json:"status"`
    Size	int	`json:"size"`
}

func resourceVolume() *schema.Resource {
    return &schema.Resource{
        Create: resourceVolumeCreate,
        Read:   resourceVolumeRead,
        Update:	resourceVolumeUpdate,
        Delete: resourceVolumeDelete,

        Timeouts: &schema.ResourceTimeout{
            Create: schema.DefaultTimeout(10 * time.Minute),
            Update: schema.DefaultTimeout(10 * time.Minute),
            Delete: schema.DefaultTimeout(10 * time.Minute),
        },

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
                ForceNew:	true,
            },

            "desc": {
                Type:		schema.TypeString,
                Optional:	true,
            },

            "is_attached": {
                Type:		schema.TypeBool,
                Computed:	true,
            },

            "is_bootable": {
                Type:		schema.TypeBool,
                Computed:	true,
                ForceNew:	true,
            },

            "is_public": {
                Type:		schema.TypeBool,
                Computed:	true,   
                ForceNew:	true,
            },

            "mountpoint": {
                Type:		schema.TypeList,
                Computed:	true,
                Elem: &schema.Schema{
                    Type:	schema.TypeString,
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

            "project": {     
                Type:		schema.TypeString,
                Optional:	true,
                Computed:	true,
                ForceNew:	true,
            },

            "size": {
                Type:		schema.TypeInt,
                Optional:	true,
                Computed:	true,
            },

            "snapshot_list": {
                Type:		schema.TypeList,
                Computed:       true,
                Elem: &schema.Schema{
                    Type:       schema.TypeString,
                },
            },

            "src_snapshot": {
                Type:		schema.TypeString,
                Optional:	true,
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

            "volume_type": {
                Type:		schema.TypeString,
                Optional:	true,
                Computed:	true,
                ForceNew:	true,
            },

            "volume_uuid": {
                Type:		schema.TypeString,
                Computed:	true,
                ForceNew:	true,
            },
        },
    }
}

func resourceVolumeCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    name := d.Get("name").(string)
    platform := d.Get("platform").(string)
    desc := d.Get("desc").(string)
    volumeType := d.Get("volume_type").(string)
    srcSnapshot := d.Get("src_snapshot").(string)
    size := d.Get("size").(int)
    project := d.Get("project").(string)
    body := VolumeCreateBody {
        Desc:		desc,
        Name:		name,
        VolumeType:	volumeType,
    }

    if srcSnapshot != "" {
        body.SrcSnapshot = srcSnapshot
    } else {
        body.Project = project
        body.Size = size
    }

    resourcePath := fmt.Sprintf("api/v3/%s/volumes/", platform)
    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    response, err := config.doNormalRequest(platform, resourcePath, "POST", buf)

    if err != nil {
        return fmt.Errorf("Error creating twcc_volume %s on %s: %v", name, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    volumeID := int(data["id"].(float64))
    d.SetId(fmt.Sprintf("%d", volumeID))

    newPath := fmt.Sprintf("%s/%d/", resourcePath, volumeID)
    stateConf := &resource.StateChangeConf{
        Pending:    []string{"CREATING", "DOWNLOADING",},
        Target:     []string{"AVAILABLE", "ERROR"},
        Refresh:    volumeStateRefreshFunc(config, platform, newPath),
        Timeout:    d.Timeout(schema.TimeoutCreate),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf(
            "Error waiting for twcc_volume %d to become AVAILABLE: %v", volumeID, err)
    }

    d.Set("name", name)
    d.Set("platform", platform)
    return resourceVolumeRead(d, meta)
}

func resourceVolumeRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    volumeID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/volumes/%s/", platform, volumeID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to retrieve volume %s on %s: %v", volumeID, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve volume json data: %v", err)
    }

    log.Printf("[DEBUG] Retrieved twcc_volume %s", d.Id())
    if host, ok := data["attached_host"].(interface{}); ok {
        hostInfo := flattenVolumeHostInfo(host)
        d.Set("attached_host", hostInfo)
    } else {
        var emptyHost interface{}
        d.Set("attached_host", emptyHost)
    }

    d.Set("create_time", data["create_time"])
    d.Set("desc", data["desc"])
    d.Set("is_attached", data["is_attached"])
    d.Set("is_bootable", data["is_bootable"])
    d.Set("is_public", data["is_public"])
    d.Set("mountpoint", data["mountpoint"])
    projectInfo := data["project"].(map[string]interface{})
    projectID := int(projectInfo["id"].(float64))
    d.Set("project", fmt.Sprintf("%d", projectID))
    d.Set("size", data["size"])
    snapshotInfo := flattenVolumeSnapshotsInfo(data["snapshot_list"].([]interface{}))
    d.Set("snapshot_list", snapshotInfo)
    d.Set("src_snapshot", data["src_snapshot"])
    d.Set("status", data["status"])
    d.Set("status_reason", data["status_reason"])
    d.Set("user", data["user"].(map[string]interface{}))
    d.Set("volume_type", data["volume_type"])
    d.Set("volume_uuid", data["volume_uuid"])
    return nil
}

func resourceVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    size_change := d.HasChange("size")
    if size_change {
        _, newSize := d.GetChange("size")
        body := VolumeUpdateBody {
            Size:	newSize.(int),
            Status:	"extend",
        }

        volumeID := d.Id()
        platform := d.Get("platform").(string)
        resourcePath := fmt.Sprintf("api/v3/%s/volumes/%s/action/", platform, volumeID)

        buf := new(bytes.Buffer)
        json.NewEncoder(buf).Encode(body)
        _, err := config.doNormalRequest(platform, resourcePath, "PUT", buf)

        if err != nil {
            return fmt.Errorf("Error resizing twcc_volume %s on %v: %s", volumeID, platform, err)
        }

        newPath := fmt.Sprintf("api/v3/%s/volumes/%s/", platform, volumeID)
        stateConf := &resource.StateChangeConf{
            Pending:    []string{"EXTENDING",},
            Target:     []string{"AVAILABLE", "ERROR", "IN-USE"},
            Refresh:    volumeStateRefreshFunc(config, platform, newPath),
            Timeout:    d.Timeout(schema.TimeoutUpdate),
            Delay:      10 * time.Second,
        }

        _, err = stateConf.WaitForState()
        if err != nil {
            return fmt.Errorf(
                "Error waiting for twcc_volume %s to become AVAILABLE: %v", volumeID, err)
        }
    }

    return resourceVolumeRead(d, meta)
}

func resourceVolumeDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    platform := d.Get("platform").(string)
    volumeID := d.Id()
    resourcePath := fmt.Sprintf("api/v3/%s/volumes/%s/", platform, volumeID)
    _, err := config.doNormalRequest(platform, resourcePath, "DELETE", nil)

    if err != nil {
        return fmt.Errorf("Unable to delete volume %s: on %s %v", volumeID, platform, err)
    }

    stateConf := &resource.StateChangeConf{
        Pending:    []string{"DELETING",},
        Target:     []string{"DELETED", "ERROR"},
        Refresh:    volumeStateRefreshForDeletedFunc(config, platform, resourcePath),
        Timeout:    d.Timeout(schema.TimeoutDelete),
        Delay:      10 * time.Second,
        MinTimeout: 3 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf( 
            "Error waiting for twcc_volume %s to become DELETED: %v", volumeID, err)
    }

    d.SetId("")

    return nil
}
