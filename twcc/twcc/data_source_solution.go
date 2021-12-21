package twcc

import (
    "fmt"
    "log"
    "encoding/json"
    "strings"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSolution() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceSolutionRead,

        Schema: map[string]*schema.Schema{
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

            "category": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "desc": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "create_time": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "is_public": {
                Type:		schema.TypeBool,
                Computed:	true,
            },

            "is_tenant_admin_only": {
                Type:		schema.TypeBool,
                Computed:	true,
            },
        },
    }
}

// dataSourceSolutionRead performs the solution lookup.
func dataSourceSolutionRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    name := d.Get("name").(string)
    params := []string{
        fmt.Sprintf("name=%s", name),
        fmt.Sprintf("project=%s", d.Get("project").(string)),
    }
    if category := d.Get("category"); category != "" {
        params = append(params, fmt.Sprintf("category=%s", category.(string)))
    }
    resourcePath := fmt.Sprintf("api/v3/solutions/?%s", strings.Join(params, "&"))
    response, err := config.doNormalRequest("goc", resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to list solutions: %v", err)
    }

    var data []map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    for _, solution := range data {
        if solution["name"] == name {
            return dataSourceSolutionAttributes(d, solution)
        }
    }

    return fmt.Errorf("Unable to retrieve solution %s: %v", name, err)
}

// dataSourceSolutionAttributes populates the fields of a solution data source.
func dataSourceSolutionAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    solution_id := int(data["id"].(float64))
    log.Printf("[DEBUG] Retrieved twcc_solution: %d", solution_id)

    d.SetId(fmt.Sprintf("%d", solution_id))
    d.Set("create_time", data["create_time"])
    d.Set("desc", data["desc"])
    d.Set("category", data["category"])
    d.Set("is_public", data["is_public"])
    d.Set("is_tenant_admin_only", data["is_tenant_admin_only"])

    return nil
}
