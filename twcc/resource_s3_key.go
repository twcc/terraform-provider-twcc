package twcc

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type KeyCreateBody struct {
    Name	string	`json:"name"`
}

func resourceS3Key() *schema.Resource {
    return &schema.Resource{
        Create: resourceS3KeyCreate,
        Read:   resourceS3KeyRead,
        Delete: resourceS3KeyDelete,

        Schema: map[string]*schema.Schema{
            "access_key": {
                Type:		schema.TypeString,
                Computed:	true,
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

            "project": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "secret_key": {
                Type:		schema.TypeString,
                Computed:	true,
            },
        },
    }
}

func resourceS3KeyCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    name := d.Get("name").(string)
    platform := d.Get("platform").(string)
    project := d.Get("project").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/projects/%s/key/", platform, project)

    body := KeyCreateBody {
        Name:	name,
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    response, err := config.doNormalRequest(platform, resourcePath, "POST", buf)

    if err != nil {
        return fmt.Errorf("Error creating twcc_s3_key %s on %s: %v", name, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    d.SetId(fmt.Sprintf("%s-%s", project, name))
    d.Set("name", name)
    d.Set("platform", platform)
    d.Set("project", project)
    return resourceS3KeyRead(d, meta)
}

func resourceS3KeyRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    projectID := d.Get("project").(string)
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/projects/%s/key/", platform, projectID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)
 
    if err != nil {
        return fmt.Errorf("Unable to retrieve s3 key on %s: %v", platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve s3 key json data: %v", err)
    }

    log.Printf("[DEBUG] Retrieved s3_key %s", d.Id())
    privateS3List := data["private"].([]interface{})
    for _, v := range privateS3List {
        key := v.(map[string]interface{})
        name := key["name"].(string)
        if name == d.Get("name").(string) {
            d.Set("access_key", key["access_key"])
            d.Set("secret_key", key["secret_key"])
            break
        }
    }
    return nil
}

func resourceS3KeyDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    name := d.Get("name").(string)
    projectID := d.Get("project").(string)
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/projects/%s/key/", platform, projectID)
    body := KeyCreateBody {
        Name:   name,
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)

    _, err := config.doNormalRequest(platform, resourcePath, "DELETE", buf)

    if err != nil {
        return fmt.Errorf("Unable to delete s3 key on %s: %v", platform, err)
    }

    d.SetId("")

    return nil
}
