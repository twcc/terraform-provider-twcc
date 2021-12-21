package twcc

import (
    "fmt"
    "log"
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceProject() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceProjectRead,

        Schema: map[string]*schema.Schema{
            "name": {
                Type:		schema.TypeString,
                Required:	true,
            },

            "platform": {
                Type:           schema.TypeString,
                Required:       true,
            },
        },
    }
}

// dataSourceProjectRead performs the project lookup.
func dataSourceProjectRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    name := d.Get("name").(string)
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/projects/", platform)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to list projects: %v", err)
    }

    var data []map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    for i := range data {
        if data[i]["name"] == name {
            return dataSourceProjectAttributes(d, data[i])
        }
    }

    return fmt.Errorf("Unable to retrieve project %s: %v", name, err)
}

// dataSourceProjectAttributes populates the fields of a project data source.
func dataSourceProjectAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    project_id := int(data["id"].(float64))
    log.Printf("[DEBUG] Retrieved twcc_project: %d", project_id)

    d.SetId(fmt.Sprintf("%d", project_id))
    d.Set("name", data["name"])
    d.Set("platform", data["platform"])

    return nil
}
