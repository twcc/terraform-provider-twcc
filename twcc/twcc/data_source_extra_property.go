package twcc

import (
    "fmt"
    "log"
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceExtraProperty() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceExtraPropertyRead,

        Schema: map[string]*schema.Schema{
            "project": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "solution": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "platform": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "extra_property": {
                Type:		schema.TypeString,
                Computed:	true,
            },
        },
    }
}

// dataSourceExtraPropertyRead performs the extra property lookup.
func dataSourceExtraPropertyRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    platform := d.Get("platform").(string)
    project := d.Get("project").(string)
    solution := d.Get("solution").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/projects/%s/solutions/%s/", platform, project, solution)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to retrieve extra property: %v", err)
    }

    var data map[string]interface{}
    if err = json.Unmarshal([]byte(response), &data); err != nil {
        return err
    }

    // remove volume-size & volume-type
    delete(data["site_extra_prop"].(map[string]interface{}), "volume-size")
    delete(data["site_extra_prop"].(map[string]interface{}), "volume-type")
    return dataSourceExtraPropertyAttributes(d, data)
}

// dataSourceExtraPropertyAttributes populates the fields of a extra property data source.
func dataSourceExtraPropertyAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    solution_id := int(data["id"].(float64))
    log.Printf("[DEBUG] Retrieved twcc_extra_property: %d", solution_id)

    json_data, err := json.Marshal(data["site_extra_prop"])
    if err != nil {
        return err
    }
    d.SetId(fmt.Sprintf("%d", solution_id))
    d.Set("extra_property", string(json_data))

    return nil
}
