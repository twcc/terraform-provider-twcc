package twcc

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type SecurityGroupRuleCreateBody struct {
    Direction	string	`json:"direction,omitempty"`
    Protocol	string	`json:"protocol,omitempty"`
    RemoteIPPrefix	string	`json:"remote_ip_prefix,omitempty"`
    PortRangeMin	int	`json:"port_range_min,omitempty"`
    PortRangeMax	int	`json:"port_range_max,omitempty"`
    Project	string	`json:"project"`
}

func resourceSecurityGroupRule() *schema.Resource {
    return &schema.Resource{
        Create: resourceSecurityGroupRuleCreate,
        Read:   resourceSecurityGroupRuleRead,
        Delete: resourceSecurityGroupRuleDelete,

        Schema: map[string]*schema.Schema{
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

            "security_group": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "direction": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "protocol": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "remote_ip_prefix": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "port_range_min": {
                Type:		schema.TypeInt,
                Optional:	true,
                ForceNew:	true,
            },

            "port_range_max": {
                Type:		schema.TypeInt,
                Optional:	true,
                ForceNew:	true,
            },

            "ethertype": {
                Type:		schema.TypeString,
                Computed:	true,
            },
        },
    }
}

func resourceSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    platform := d.Get("platform").(string)
    projectID := d.Get("project").(string)
    securityGroupID := d.Get("security_group").(string)
    direction := d.Get("direction").(string)
    protocol := d.Get("protocol").(string)
    remote_ip_prefix := d.Get("remote_ip_prefix").(string)
    port_range_min := d.Get("port_range_min").(int)
    port_range_max := d.Get("port_range_max").(int)
    if port_range_min == 0 && port_range_max != 0 {
        port_range_min = port_range_max
    } else if port_range_max == 0 && port_range_min != 0 {
        port_range_max = port_range_min
    }
    resourcePath := fmt.Sprintf("api/v3/%s/security_groups/%s/", platform, securityGroupID)

    body := SecurityGroupRuleCreateBody {
        Direction:		direction,
        Protocol:		protocol,
        RemoteIPPrefix:		remote_ip_prefix,
        PortRangeMin:		port_range_min,
        PortRangeMax:		port_range_max,
        Project:	projectID,
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    _, err := config.doNormalRequest(platform, resourcePath, "PATCH", buf)

    if err != nil {
        return fmt.Errorf(
            "Error creating twcc_security_group_rule from %s on %s: %v",
            securityGroupID,
            platform,
            err,
        )
    }

    return resourceSecurityGroupRuleRead(d, meta)
}

func resourceSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    platform := d.Get("platform").(string)
    projectID := d.Get("project").(string)
    securityGroupID := d.Get("security_group").(string)
    resourcePath := fmt.Sprintf(
        "api/v3/%s/security_groups/?project=%s&sg=%s",
        platform,
        projectID,
        securityGroupID,
    )

    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf(
            "Unable to retrieve security group %s on %s: %v", securityGroupID, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve security group json data: %v", err)
    }

    if security_group_rules, ok := data["security_group_rules"].([]interface{}); ok {
        for _, security_group_rule := range security_group_rules {
            sg_rule := security_group_rule.(map[string]interface{})
            if foundSecurityGroupRule(sg_rule, d) {
                securityGroupRuleID := sg_rule["id"].(string)
                log.Printf("[DEBUG] Retrieved twcc_security_group_rule %s", securityGroupRuleID)
                d.SetId(securityGroupRuleID)
                d.Set("direction", sg_rule["direction"].(string))
                d.Set("protocol", sg_rule["protocol"].(string))
                d.Set("remote_ip_prefix", sg_rule["remote_ip_prefix"].(string))
                d.Set("ethertype", sg_rule["ethertype"].(string))
                if port_range_min, ok := sg_rule["port_range_min"].(float64); ok {
                    d.Set("port_range_min", int(port_range_min))
                }
                if port_range_max, ok := sg_rule["port_range_max"].(float64); ok {
                    d.Set("port_range_max", int(port_range_max))
                }
                return nil
            }
        }
    }
    return fmt.Errorf("Unable to retrieve security group rule from %s on %s", securityGroupID, platform)
}

func resourceSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    platform := d.Get("platform").(string)
    projectID := d.Get("project").(string)
    securityGroupRuleID := d.Id()
    resourcePath := fmt.Sprintf(
        "api/v3/%s/security_group_rules/%s/?project=%s",
        platform,
        securityGroupRuleID,
        projectID,
    )

    _, err := config.doNormalRequest(platform, resourcePath, "DELETE", nil)

    if err != nil {
        return fmt.Errorf(
            "Unable to delete security group rule %s on %s: %v",
            securityGroupRuleID,
            platform,
            err,
        )
    }

    d.SetId("")
    return nil
}
