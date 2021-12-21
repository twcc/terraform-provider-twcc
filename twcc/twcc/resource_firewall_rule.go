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

type FirewallRuleCreateBody struct {
    Action			string	`json:"action,omitempty"`
    DestinationIPAddress	string	`json:"destination_ip_address,omitempty"`
    DestinationPort		string	`json:"destination_port,omitempty"`
    Name			string	`json:"name"`
    Project			string	`json:"project"`
    Protocol			string	`json:"protocol,omitempty"`
    SourceIPAddress		string	`json:"source_ip_address,omitempty"`
    SourcePort			string	`json:"source_port,omitempty"`
}

type FirewallRuleUpdateBody struct {
    Action			string	`json:"action,omitempty"`
    DestinationIPAddress	string	`json:"destination_ip_address,omitempty"`
    DestinationPort		string	`json:"destination_port,omitempty"`
    Protocol			string	`json:"protocol,omitempty"`
    SourceIPAddress		string	`json:"source_ip_address,omitempty"`
    SourcePort			string	`json:"source_port,omitempty"`
}

func resourceFirewallRule() *schema.Resource {
    return &schema.Resource{
        Create: resourceFirewallRuleCreate,
        Read:   resourceFirewallRuleRead,
        Update: resourceFirewallRuleUpdate,
        Delete: resourceFirewallRuleDelete,

        Timeouts: &schema.ResourceTimeout{
            Delete: schema.DefaultTimeout(10 * time.Minute),
        },

        Schema: map[string]*schema.Schema{
            "action": {
                Type:		schema.TypeString,
                Optional:	true,
                Computed:	true,
            },

            "create_time": {
                Type:		schema.TypeString,
                Computed:	true,
                ForceNew:	true,
            },

            "destination_ip_address": {
                Type:		schema.TypeString,
                Optional:	true,
            },

            "destination_port": {
                Type:		schema.TypeString,
                Optional:	true,
            },

            "ip_version": {
                Type:		schema.TypeInt,
                Computed:	true,
                ForceNew:	true,
            },

            "name": {
                Type:		schema.TypeString,
                Required:	true,
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

            "protocol": {
                Type:		schema.TypeString,
                Optional:	true,
                Default:	"tcp",
            },

            "source_ip_address": {
                Type:		schema.TypeString,
                Optional:	true,
            },

            "source_port": {
                Type:		schema.TypeString,
                Optional:	true,
            },

        },
    }
}

func resourceFirewallRuleCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    action := d.Get("action").(string)
    destinationIPAddress := d.Get("destination_ip_address").(string)
    destinationPort :=  d.Get("destination_port").(string)
    name := d.Get("name").(string)
    platform := d.Get("platform").(string)
    project := d.Get("project").(string)
    protocol := d.Get("protocol").(string)
    sourceIPAddress := d.Get("source_ip_address").(string)
    sourcePort :=  d.Get("source_port").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/firewall_rules/", platform)

    body := FirewallRuleCreateBody {
        Action:			action,
        DestinationIPAddress:	destinationIPAddress,
        DestinationPort:	destinationPort,
        Name:			name,
        Project:		project,
        Protocol:		protocol,
        SourceIPAddress:	sourceIPAddress,
        SourcePort:	 	sourcePort,
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    response, err := config.doNormalRequest(platform, resourcePath, "POST", buf)

    if err != nil {
        return fmt.Errorf("Error creating twcc_firewall_rule %s on %s: %v", name, platform, err)
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
    return resourceFirewallRuleRead(d, meta)
}

func resourceFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    ruleID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/firewall_rules/%s/", platform, ruleID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to retrieve firewall rule %s on %s: %v", ruleID, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve firewall rule json data: %v", err)
    }

    log.Printf("[DEBUG] Retrieved twcc_firewall_rule %d", d.Id())
    d.Set("action", data["action"])
    d.Set("create_time", data["create_time"])
    d.Set("destination_ip_address", data["destination_ip_address"])
    d.Set("destination_port", data["destination_port"])
    d.Set("protocol", data["protocol"])
    d.Set("source_ip_address", data["source_ip_address"])
    d.Set("source_port", data["source_port"])
    return nil
}

func resourceFirewallRuleUpdate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    if d.HasChange("action") || d.HasChange("destination_ip_address") ||
            d.HasChange("destination_port") || d.HasChange("protocol") ||
            d.HasChange("source_ip_address") || d.HasChange("source_port") {
        var body FirewallRuleUpdateBody
        if d.HasChange("action") {
            _, newAction := d.GetChange("action")
            body.Action = newAction.(string)
        }

        if d.HasChange("destination_ip_address") {
            _, newDestinationIPAddress := d.GetChange("destination_ip_address")
            body.DestinationIPAddress = newDestinationIPAddress.(string)
        }

        if d.HasChange("destination_port") {
            _, newDestinationPort := d.GetChange("destination_port")
            body.DestinationPort = newDestinationPort.(string)
        }

        if d.HasChange("protocol") {
            _, newProtocol := d.GetChange("protocol")
            body.Protocol = newProtocol.(string)
        }

        if d.HasChange("source_ip_address") {
            _, newSourceIPAddress := d.GetChange("source_ip_address")
            body.SourceIPAddress = newSourceIPAddress.(string)
        }

        if d.HasChange("source_port") {
            _, newSourcePort := d.GetChange("source_port")
            body.SourcePort = newSourcePort.(string)
        }

        ruleID := d.Id()
        platform := d.Get("platform").(string)
        resourcePath := fmt.Sprintf("api/v3/%s/firewall_rules/%s/", platform, ruleID)

        buf := new(bytes.Buffer)
        json.NewEncoder(buf).Encode(body)
        _, err := config.doNormalRequest(platform, resourcePath, "PATCH", buf)

        if err != nil {
            return fmt.Errorf("Error updating twcc_firewall_rule %s on %s: %v", ruleID, platform, err)
        }
    }

    return resourceFirewallRuleRead(d, meta)
}

func resourceFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    ruleID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/firewall_rules/%s/", platform, ruleID)
    _, err := config.doNormalRequest(platform, resourcePath, "DELETE", nil)

    if err != nil {
        return fmt.Errorf("Unable to delete firewall rule %s: on %s %v", ruleID, platform, err)
    }

    stateConf := &resource.StateChangeConf{
        Pending:    []string{"DELETING"},
        Target:     []string{"DELETED"},
        Refresh:    firewallRuleStateRefreshForDeletedFunc(config, platform, resourcePath),
        Timeout:    d.Timeout(schema.TimeoutCreate),
        Delay:      10 * time.Second,
        MinTimeout: 3 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf(
            "Error waiting for twcc_firewall_rule %s to become DELETED: %v", ruleID, err)
    }

    d.SetId("")

    return nil
}
