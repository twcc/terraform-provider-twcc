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

type AutoScalingPolicyCreateBody struct {
    Description		string	`json:"description,omitempty"`
    MeterName		string	`json:"meter_name"`
    Name		string	`json:"name"`
    Project		string	`json:"project"`
    ScaledownThreshold	int	`json:"scaledown_threshold"`
    ScaleMaxSize	int	`json:"scale_max_size"`
    ScaleupThreshold	int	`json:"scaleup_threshold"`
}

func resourceAutoScalingPolicy() *schema.Resource {
    return &schema.Resource{
        Create: resourceAutoScalingPolicyCreate,
        Read:   resourceAutoScalingPolicyRead,
        Delete: resourceAutoScalingPolicyDelete,

        Timeouts: &schema.ResourceTimeout{
            Delete: schema.DefaultTimeout(15 * time.Minute),
        },

        Schema: map[string]*schema.Schema{
            "description": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "meter_name": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
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

            "scale_max_size": {
                Type:		schema.TypeInt,
                Required:	true,
                ForceNew:	true,
            },

            "scaledown_threshold": {
                Type:		schema.TypeInt,
                Required:	true,
                ForceNew:	true,
            },

            "scaleup_threshold": {
                Type:		schema.TypeInt,
                Required:	true,
                ForceNew:	true,
            },

            "user": {       
                Type:		schema.TypeMap,
                Computed:	true,
                ForceNew:	true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
        },
    }
}

func resourceAutoScalingPolicyCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    description := d.Get("description").(string)
    meterName := d.Get("meter_name").(string)
    name := d.Get("name").(string)
    platform := d.Get("platform").(string)
    project := d.Get("project").(string)
    scaleMaxSize := d.Get("scale_max_size").(int)
    scaledownThreshold := d.Get("scaledown_threshold").(int)
    scaleupThreshold := d.Get("scaleup_threshold").(int)
    resourcePath := fmt.Sprintf("api/v3/%s/auto_scaling_policies/", platform)

    body := AutoScalingPolicyCreateBody {
        Description:		description,
        MeterName:		meterName,
        Name:			name,
        Project:		project,
    	ScaledownThreshold:	scaledownThreshold,
        ScaleMaxSize:		scaleMaxSize,
        ScaleupThreshold:	scaleupThreshold,
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    response, err := config.doNormalRequest(platform, resourcePath, "POST", buf)

    if err != nil {
        return fmt.Errorf("Error creating twcc_auto_scaling_policy %s on %s: %v", name, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    d.SetId(fmt.Sprintf("%d", int(data["id"].(float64))))
    d.Set("name", name)
    d.Set("platform", platform)
    d.Set("project", project)
    return resourceAutoScalingPolicyRead(d, meta)
}

func resourceAutoScalingPolicyRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    policyID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/auto_scaling_policies/%s/", platform, policyID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)
 
    if err != nil {
        return fmt.Errorf("Unable to retrieve auto scaling policy %s on %s: %v", policyID, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve auto scaling policy json data: %v", err)
    }

    log.Printf("[DEBUG] Retrieved twcc_auto_scaling_policy %s", d.Id())
    d.Set("description", data["description"])
    d.Set("meter_name", data["meter_name"])
    d.Set("scale_max_size", data["scale_max_size"])
    d.Set("scaledown_threshold", data["scaledown_threshold"])
    d.Set("scaleup_threshold", data["scaleup_threshold"])
    d.Set("user", data["user"])
    return nil
}

func resourceAutoScalingPolicyDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    policyID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/auto_scaling_policies/%s/", platform, policyID)
    _, err := config.doNormalRequest(platform, resourcePath, "DELETE", nil)

    if err != nil {
        return fmt.Errorf("Unable to delete auto scaling policy %s: on %s %v", policyID, platform, err)
    }

    stateConf := &resource.StateChangeConf{
        Pending:    []string{"DELETING"},
        Target:     []string{"DELETED"},
        Refresh:    policyStateRefreshForDeletedFunc(config, platform, resourcePath),
        Timeout:    d.Timeout(schema.TimeoutDelete),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf( 
            "Error waiting for twcc_auto_scaling_policy %s to become Deleted: %v", policyID, err)
    }

    d.SetId("")

    return nil
}
