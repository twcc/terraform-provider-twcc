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

type ServerImageCreateBody struct {
    Desc	string	`json:"desc,omitempty"`
    Name	string	`json:"name"`
    OS		string	`json:"os"`
    OSVersion	string	`json:"os_version"`
}

func resourceVCSImage() *schema.Resource {
    return &schema.Resource{
        Create: resourceVCSImageCreate,
        Read:   resourceVCSImageRead,
        Delete: resourceVCSImageDelete,

        Timeouts: &schema.ResourceTimeout{
            Create: schema.DefaultTimeout(30 * time.Minute),
            Delete: schema.DefaultTimeout(15 * time.Minute),
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

            "is_enabled": {
                Type:		schema.TypeBool,
                Computed:	true,
            },

            "is_public": {
                Type:		schema.TypeBool,
                Computed:	true,
            },

            "name": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "os": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "os_version": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "platform": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "ref_img_id": {
                Type:		schema.TypeString,
                Computed:	true,
                ForceNew:	true,
            },

            "server": {
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
        },
    }
}

func resourceVCSImageCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    desc := d.Get("desc").(string)
    name := d.Get("name").(string)
    os := d.Get("os").(string)
    osVersion := d.Get("os_version").(string)
    platform := d.Get("platform").(string)
    server := d.Get("server").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/images/%s/save/", platform, server)

    body := ServerImageCreateBody {
        Desc:		desc,
        Name:		name,
        OS:		os,
        OSVersion:	osVersion,
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    response, err := config.doNormalRequest(platform, resourcePath, "PUT", buf)

    if err != nil {
        return fmt.Errorf(
            "Error creating twcc_vcs_image %s from %s on %s: %v",
            name,
            server,
            platform,
            err,
        )
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    imageID := int(data["id"].(float64))
    d.SetId(fmt.Sprintf("%d", imageID))

    newPath := fmt.Sprintf("api/v3/%s/images/%d/", platform, imageID)
    stateConf := &resource.StateChangeConf{
        Pending:    []string{"QUEUED", "SAVING"},
        Target:     []string{"ACTIVE", "ERROR"},
        Refresh:    imageStateRefreshFunc(config, platform, newPath),
        Timeout:    d.Timeout(schema.TimeoutCreate),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf("Error waiting for twcc_vcs_image %d to become ACTIVE: %v", imageID, err)
    }

    d.Set("desc", desc)
    d.Set("name", name)
    d.Set("os", os)
    d.Set("os_version", osVersion)
    d.Set("platform", platform)
    d.Set("server", server) 
    return resourceVCSImageRead(d, meta)
}

func resourceVCSImageRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    imageID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/images/%s/", platform, imageID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf(
            "Unable to retrieve image %s on %s: %v", imageID, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve image json data: %s", err)
    }

    log.Printf("[DEBUG] Retrieved twcc_vcs_image %s", imageID)
    d.Set("create_time", data["create_time"].(string))
    d.Set("is_enabled", data["is_enabled"].(bool))
    d.Set("is_public", data["is_public"].(bool))
    d.Set("ref_img_id", data["ref_img_id"].(string))
    d.Set("status", data["status"].(string))
    d.Set("status_reason", data["status_reason"].(string))
    return nil
}

func resourceVCSImageDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    imageID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/images/%s/", platform, imageID)
    _, err := config.doNormalRequest(platform, resourcePath, "DELETE", nil)

    if err != nil {
        return fmt.Errorf("Unable to delete VCS snapshot image %s on %s: %v", imageID, platform, err)
    }

    stateConf := &resource.StateChangeConf{
        Pending:    []string{"ACTIVE"},
        Target:     []string{"DELETED", "ERROR"},
        Refresh:    imageStateRefreshForDeletedFunc(config, platform, resourcePath),
        Timeout:    d.Timeout(schema.TimeoutDelete),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf( 
            "Error waiting for twcc_vcs_image to become DELETED: %v", err)
    }

    d.SetId("")

    return nil
}
