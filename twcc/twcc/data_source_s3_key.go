package twcc

import (
    "fmt"
    "log"
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceS3Key() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceS3KeyRead,

        Schema: map[string]*schema.Schema{
            "access_key": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "is_public": {
                Type:		schema.TypeBool,
                Computed:       true,
                ForceNew:       true,
            },

            "name": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "platform": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "project": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "secret_key": {
                Type:		schema.TypeString,
                Computed:   true,
            },
        },
    }
}

// dataSourceS3KeyRead performs the s3 key lookup.
func dataSourceS3KeyRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    name := d.Get("name").(string)
    projectID := d.Get("project").(string)
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/projects/%s/key/", platform, projectID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)
    if err != nil {
        return fmt.Errorf("Unable to retrive project keys: %v", err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    d.Set("platform", platform)
    d.Set("project", projectID)
    if name != "" {
        d.Set("is_public", false)
        privateS3List := data["private"].([]interface{})
        for _, v := range privateS3List{
            key := v.(map[string]interface{})
            if name == key["name"].(string) {
                return dataSourceS3KeyAttributes(d, key)
            }
        }
    } else {
        d.Set("is_public", true)
        return dataSourceS3KeyAttributes(d, data["public"].(map[string]interface{}))
    }

    return fmt.Errorf("Unable to retrieve s3 key %s", name)
}

// dataSourceS3KeyAttributes populates the fields of a s3 key data source.
func dataSourceS3KeyAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    log.Printf("[DEBUG] Retrieved twcc_s3_key")
    project := d.Get("project").(string)
    if name, ok := data["name"].(string); ok {
        d.Set("name", name)
        d.SetId(fmt.Sprintf("%s-%s", project, name))
    } else {
        d.Set("name", "Public Key")
        d.SetId(fmt.Sprintf("%s-%s", project, "public"))
    }

    d.Set("access_key", data["access_key"])
    d.Set("secret_key", data["secret_key"])

    return nil
}
