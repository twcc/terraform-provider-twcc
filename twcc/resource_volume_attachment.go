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

type VolumeAttachmentCreateBody struct {
    Mountpoint	string	`json:"mountpoint,omitempty"`
    Server	string	`json:"server,omitempty"`
    Status	string	`json:"status"`
}

func resourceVolumeAttachment() *schema.Resource {
    return &schema.Resource{
        Create: resourceVolumeAttachmentCreate,
        Read:   resourceVolumeAttachmentRead,
        Delete: resourceVolumeAttachmentDelete,

        Timeouts: &schema.ResourceTimeout{
            Create: schema.DefaultTimeout(10 * time.Minute),
            Delete: schema.DefaultTimeout(10 * time.Minute),
        },

        Schema: map[string]*schema.Schema{
            "mountpoint": {
                Type:		schema.TypeString,
                Optional:	true,
                Computed:	true,
            },

            "platform": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "server": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "volume": { 
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },
        },
    }
}

func resourceVolumeAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    mountpoint := d.Get("mountpoint").(string)
    platform := d.Get("platform").(string)
    server := d.Get("server").(string)
    volume := d.Get("volume").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/volumes/%s/action/", platform, volume)

    body := VolumeAttachmentCreateBody {
        Mountpoint:	mountpoint,
        Server:		server,
        Status:		"attach",
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    _, err := config.doNormalRequest(platform, resourcePath, "PUT", buf)

    if err != nil {
        return fmt.Errorf(
            "Error creating twcc_volume_attachment with volume %s and server %s on %s: %v",
            volume,
            server,
            platform,
            err,
        )
    }

    d.SetId(fmt.Sprintf("%s/%s", server, volume))

    newPath := fmt.Sprintf("api/v3/%s/volumes/%s/", platform, volume)
    stateConf := &resource.StateChangeConf{
        Pending:    []string{"ATTACHING"},
        Target:     []string{"IN-USE", "ERROR"},
        Refresh:    volumeStateRefreshFunc(config, platform, newPath),
        Timeout:    d.Timeout(schema.TimeoutCreate),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf(
            "Error waiting for twcc_volume_attachment with volume %s and server %s to become IN-USE: %v",
            volume,
            server,
            err,
        )
    }

    d.Set("platform", platform)
    d.Set("server", server)
    d.Set("volume", volume)
    return resourceVolumeAttachmentRead(d, meta)
}

func resourceVolumeAttachmentRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    volumeID := d.Get("volume").(string)
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/volumes/%s/", platform, volumeID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf(
            "Unable to retrieve volume %s on %s: %s", volumeID, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve volume json data: %v", err)
    }

    log.Printf("[DEBUG] Retrieved twcc_volume_attachment by volume %s", volumeID)
    mountpoints := data["mountpoint"].([]interface{})
    d.Set("mountpoint", mountpoints[0].(string))
    return nil
}

func resourceVolumeAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    serverID := d.Get("server").(string)
    volumeID := d.Get("volume").(string)
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/volumes/%s/action/", platform, volumeID)

    body := VolumeAttachmentCreateBody {
        Server:	serverID,
        Status:	"detach",
    }
    
    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    _, err := config.doNormalRequest(platform, resourcePath, "PUT", buf)

    if err != nil {
        return fmt.Errorf(
            "Unable to delete volume attachment with volume %s and server %s on %s: %v",
            volumeID,
            serverID,
            platform,
            err,
        )
    }

    newPath := fmt.Sprintf("api/v3/%s/volumes/%s/", platform, volumeID)
    stateConf := &resource.StateChangeConf{
        Pending:    []string{"DETACHING"},
        Target:     []string{"AVAILABLE", "ERROR"},
        Refresh:    volumeStateRefreshFunc(config, platform, newPath),
        Timeout:    d.Timeout(schema.TimeoutDelete),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf( 
            "Error waiting for twcc_volume_attachment to become AVAILABLE: %v", err)
    }

    d.SetId("")

    return nil
}
