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

type AutoScalingRelationCreateBody struct {
    AutoScalingPolicy	string	`json:"auto_scaling_policy"`
    Loadbalancer	string	`json:"loadbalancer,omitempty"`
    ProtocolPort	int	`json:"protocol_port,omitempty"`
    ScaledownAction	string	`json:"scaledown_action,omitempty"`
    ScaleupAction	string	`json:"scaleup_action,omitempty"`
}

func resourceAutoScalingRelation() *schema.Resource {
    return &schema.Resource{
        Create: resourceAutoScalingRelationCreate,
        Read:   resourceAutoScalingRelationRead,
        Delete: resourceAutoScalingRelationDelete,

        Timeouts: &schema.ResourceTimeout{
            Create: schema.DefaultTimeout(15 * time.Minute),
            Delete: schema.DefaultTimeout(15 * time.Minute),
        },

        Schema: map[string]*schema.Schema{
            "auto_scaling_policy": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "loadbalancer": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "platform": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "protocol_port": {     
                Type:		schema.TypeInt,
                Optional:	true,
                ForceNew:	true,
            },

            "scaledown_action": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "scaleup_action": {
                Type:		schema.TypeString,
                Optional:	true,
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

func resourceAutoScalingRelationCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    autoScalingPolicy := d.Get("auto_scaling_policy").(string)
    loadbalancer := d.Get("loadbalancer").(string)
    platform := d.Get("platform").(string)
    protocolPort := d.Get("protocol_port").(int)
    scaledownAction := d.Get("scaledown_action").(string)
    scaleupAction := d.Get("scaleup_action").(string)
    server := d.Get("server").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/servers/%s/auto_scaling_policy/", platform, server)

    body := AutoScalingRelationCreateBody {
        AutoScalingPolicy:	autoScalingPolicy,
        Loadbalancer:		loadbalancer,
        ProtocolPort:		protocolPort,
        ScaledownAction:	scaledownAction,
        ScaleupAction:		scaleupAction,
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    _, err := config.doNormalRequest(platform, resourcePath, "POST", buf)

    if err != nil {
        return fmt.Errorf(
            "Error creating twcc_auto_scaling_relation with policy %s and server %s on %s: %v",
            autoScalingPolicy,
            server,
            platform,
            err,
        )
    }

    d.SetId(fmt.Sprintf("%s/%s", server, autoScalingPolicy))

    newPath := fmt.Sprintf("api/v3/%s/servers/%s/", platform, server)
    stateConf := &resource.StateChangeConf{
        Pending:    []string{"ASSOCIATING"},
        Target:     []string{"ASSOCIATED", "ERROR"},
        Refresh:    relationStateRefreshFunc(config, platform, newPath),
        Timeout:    d.Timeout(schema.TimeoutCreate),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf(
            "Error waiting for twcc_auto_scaling_relation with policy %s and server %s to become Ready: %v",
            autoScalingPolicy,
            server,
            err,
        )
    }

    d.Set("auto_scaling_policy", autoScalingPolicy)
    d.Set("loadbalancer", loadbalancer)
    d.Set("platform", platform)
    d.Set("protocol_port", protocolPort)
    d.Set("scaledown_action", scaledownAction)
    d.Set("scaleup_action", scaleupAction)
    d.Set("server", server) 
    return resourceAutoScalingRelationRead(d, meta)
}

func resourceAutoScalingRelationRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    serverID := d.Get("server").(string)
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/servers/%s/", platform, serverID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf(
            "Unable to retrieve server %s on %s: %v", serverID, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve server json data: %v", err)
    }

    log.Printf("[DEBUG] Retrieved twcc_auto_scaling_relation by server %s", serverID)
    info := data["auto_scaling_policy"].(map[string]interface{})
    d.Set("status", info["status"].(string))
    d.Set("status_reason", info["status_reason"].(string))
    return nil
}

func resourceAutoScalingRelationDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    serverID := d.Get("server").(string)
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/servers/%s/auto_scaling_policy/", platform, serverID)
    _, err := config.doNormalRequest(platform, resourcePath, "DELETE", nil)

    if err != nil {
        return fmt.Errorf(
            "Unable to delete auto scaling relation with policy %s and server %s on %s: %v",
            d.Get("auto_scaling_policy").(string),
            serverID,
            platform,
            err,
        )
    }

    newPath := fmt.Sprintf("api/v3/%s/servers/%s/", platform, serverID)
    stateConf := &resource.StateChangeConf{
        Pending:    []string{"DISASSOCIATING"},
        Target:     []string{"DELETED", "ERROR"},
        Refresh:    relationStateRefreshForDeletedFunc(config, platform, newPath),
        Timeout:    d.Timeout(schema.TimeoutDelete),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf( 
            "Error waiting for twcc_auto_scaling_relation to become Deleted: %v", err)
    }

    d.SetId("")

    return nil
}
