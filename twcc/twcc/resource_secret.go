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

type SecretCreateBody struct {
    Desc	string	`json:"desc,omitempty"`
    ExpireTime	string	`json:"expire_time,omitempty"`
    Name	string	`json:"name"`
    Payload	string	`json:"payload"`
    Project	string	`json:"project"`
}

func resourceSecret() *schema.Resource {
    return &schema.Resource{
        Create: resourceSecretCreate,
        Read:   resourceSecretRead,
        Delete: resourceSecretDelete,

        Timeouts: &schema.ResourceTimeout{
            Delete: schema.DefaultTimeout(5 * time.Minute),
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

            "expire_time": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "name": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:       true,
            },

            "payload": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "platform": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:       true,
            },

            "project": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "status": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "user": {
                Type:	schema.TypeMap,
                Computed:	true,
                ForceNew:	true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
        },
    }
}

func resourceSecretCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    desc := d.Get("desc").(string)
    expireTime := d.Get("expire_time").(string)
    name := d.Get("name").(string)
    payload := d.Get("payload").(string)
    platform := d.Get("platform").(string)
    project := d.Get("project").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/secrets/", platform)

    body := SecretCreateBody {
        Desc:		desc,
        ExpireTime:	expireTime,
        Name:		name,
        Payload:	payload,
        Project:	project,
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    response, err := config.doNormalRequest(platform, resourcePath, "POST", buf)

    if err != nil {
        return fmt.Errorf("Error creating twcc_secret %s on %s: %v", name, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    secretID := int(data["id"].(float64))
    d.SetId(fmt.Sprintf("%d", secretID))
    d.Set("name", name)
    d.Set("platform", platform)
    d.Set("project", project)
    d.Set("desc", desc)
    d.Set("payload", payload)
    d.Set("expire_time", expireTime)
    return resourceSecretRead(d, meta)
}

func resourceSecretRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    secretID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/secrets/%s/", platform, secretID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)
 
    if err != nil {
        return fmt.Errorf("Unable to retrieve secret %s on %s: %v", secretID, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve secret json data: %v", err)
    }

    log.Printf("[DEBUG] Retrieved secret %s", d.Id())
    d.Set("create_time", data["create_time"])
    d.Set("status", data["status"])
    d.Set("user", data["user"].(map[string]interface{}))
    return nil
}

func resourceSecretDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    secretID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/secrets/%s/", platform, secretID)
    _, err := config.doNormalRequest(platform, resourcePath, "DELETE", nil)

    if err != nil {
        return fmt.Errorf("Unable to delete secret %s on %s: %v", secretID, platform, err)
    }

    stateConf := &resource.StateChangeConf{
        Pending:    []string{"DELETING"},
        Target:     []string{"DELETED", "ERROR"},
        Refresh:    secretStateRefreshForDeletedFunc(config, platform, resourcePath),
        Timeout:    d.Timeout(schema.TimeoutDelete),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf(
            "Error waiting for twcc_secret %s to become DELETED: %v", secretID, err)
    }

    d.SetId("")

    return nil
}
