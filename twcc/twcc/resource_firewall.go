package twcc

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "time"
    "strconv"

    "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type FirewallCreateBody struct {
    AssociateNetworks	[]int	`json:"associate_networks,omitempty"`
    Desc		string	`json:"desc,omitempty"`
    Name		string	`json:"name"`
    Project		string	`json:"project"`
    Rules		[]int	`json:"rules,omitempty"`
}

type FirewallUpdateBody struct {
    AssociateNetworks   *[]int	`json:"associate_networks,omitempty"`
    Desc                string	`json:"desc,omitempty"`
    Rules               *[]int	`json:"rules,omitempty"`
}

func resourceFirewall() *schema.Resource {
    return &schema.Resource{
        Create: resourceFirewallCreate,
        Read:   resourceFirewallRead,
        Update:	resourceFirewallUpdate,
        Delete: resourceFirewallDelete,

        Timeouts: &schema.ResourceTimeout{
            Create: schema.DefaultTimeout(10 * time.Minute),
            Update: schema.DefaultTimeout(10 * time.Minute),
            Delete: schema.DefaultTimeout(10 * time.Minute),
        },

        Schema: map[string]*schema.Schema{
            "associate_networks": {
                Type:		schema.TypeList,
                Optional:	true,
                Elem: &schema.Schema{
                    Type:	schema.TypeString,
                },
            },

            "create_time": {
                Type:		schema.TypeString,
                Computed:	true,
                ForceNew:	true,
            },

            "desc": {
                Type:		schema.TypeString,
                Optional:	true,
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

            "rules": {
                Type:		schema.TypeList,
                Optional:	true,
                Elem: &schema.Schema{
                    Type:	schema.TypeString,
                },
            },

            "status": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "status_reason": {
                Type:		schema.TypeString,
                Computed:	true,
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

func resourceFirewallCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    name := d.Get("name").(string)
    platform := d.Get("platform").(string)
    project := d.Get("project").(string)
    networks := d.Get("associate_networks").([]interface{})
    networkIDArray := make([]int, len(networks))
    for i, net := range networks {
        IDString := net.(string)
        IDInt, err := strconv.Atoi(IDString)
        if err != nil {
            return fmt.Errorf("Not correct ID format %s", IDString)
        }

        networkIDArray[i] = IDInt
    }

    rules := d.Get("rules").([]interface{})
    ruleIDArray := make([]int, len(rules))
    for i, rule := range rules {
        IDString := rule.(string)
        IDInt, err := strconv.Atoi(IDString)
        if err != nil {
            return fmt.Errorf("Not correct ID format %s", IDString)
        }

        ruleIDArray[i] = IDInt
    }


    desc := d.Get("desc").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/firewalls/", platform)

    body := FirewallCreateBody {
        AssociateNetworks:	networkIDArray,
        Desc:			desc,
        Name:			name,
        Project:		project,
        Rules:			ruleIDArray,
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    response, err := config.doNormalRequest(platform, resourcePath, "POST", buf)

    if err != nil {
        return fmt.Errorf("Error creating twcc_firewall %s on %s: %v", name, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    firewallID := int(data["id"].(float64))
    d.SetId(fmt.Sprintf("%d", firewallID))

    newPath := fmt.Sprintf("%s/%d/", resourcePath, firewallID)
    stateConf := &resource.StateChangeConf{
        Pending:    []string{"PENDING_UPDATE", "PENDING_DELETE",},
        Target:     []string{"ACTIVE", "ERROR", "INACTIVE"},
        Refresh:    firewallStateRefreshFunc(config, platform, newPath),
        Timeout:    d.Timeout(schema.TimeoutCreate),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf(
            "Error waiting for twcc_firewall %d to become ACTIVE: %v", firewallID, err)
    }

    d.Set("name", name)
    d.Set("platform", platform)
    d.Set("project", project)
    return resourceFirewallRead(d, meta)
}

func resourceFirewallRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    firewallID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/firewalls/%s/", platform, firewallID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to retrieve firewall %s on %s: %v", firewallID, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve firewall json data: %v", err)
    }

    log.Printf("[DEBUG] Retrieved twcc_firewall %s", d.Id())
    networkInfo := flattenFirewallObjectInfo(data["associate_networks"].([]interface{}))
    d.Set("associate_networks", networkInfo)
    d.Set("create_time", data["create_time"])
    d.Set("desc", data["desc"])
    ruleInfo := flattenFirewallObjectInfo(data["rules"].([]interface{}))
    d.Set("rules", ruleInfo)
    d.Set("status", data["status"])
    d.Set("status_reason", data["status_reason"])
    d.Set("user", data["user"].(map[string]interface{}))
    return nil
}

func resourceFirewallUpdate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    a_change := d.HasChange("associate_networks")
    d_change := d.HasChange("desc")
    r_change := d.HasChange("rules")
    if a_change || d_change || r_change {
        var body FirewallUpdateBody
        if a_change {
            _, newAssociateNetworks := d.GetChange("associate_networks")
            networks := newAssociateNetworks.([]interface{})
            networkIDArray := make([]int, len(networks))
            for i, net := range networks {
                IDString := net.(string)
                IDInt, err := strconv.Atoi(IDString)
                if err != nil {
                    return fmt.Errorf("Not correct ID format %s", IDString)
                }

                networkIDArray[i] = IDInt
            }

            body.AssociateNetworks = &networkIDArray
        } else if r_change {
            _, newRules := d.GetChange("rules")
            rules := newRules.([]interface{})
            ruleIDArray := make([]int, len(rules))
            for i, rule := range rules {
                IDString := rule.(string)
                IDInt, err := strconv.Atoi(IDString)
                if err != nil {
                    return fmt.Errorf("Not correct ID format %s", IDString)
                }

                ruleIDArray[i] = IDInt
            }

            body.Rules = &ruleIDArray
        }

        if d_change {
            _, newDesc := d.GetChange("desc")
            body.Desc = newDesc.(string)
        }

        firewallID := d.Id()
        platform := d.Get("platform").(string)
        resourcePath := fmt.Sprintf("api/v3/%s/firewalls/%s/", platform, firewallID)

        buf := new(bytes.Buffer)
        json.NewEncoder(buf).Encode(body)
        _, err := config.doNormalRequest(platform, resourcePath, "PATCH", buf)

        if err != nil {
            return fmt.Errorf("Error updating twcc_firewall %s on %s: %v", firewallID, platform, err)
        }

        stateConf := &resource.StateChangeConf{
            Pending:    []string{"PENDING_UPDATE", "PENDING_DELETE",},
            Target:     []string{"ACTIVE", "ERROR", "INACTIVE"},
            Refresh:    firewallStateRefreshFunc(config, platform, resourcePath),
            Timeout:    d.Timeout(schema.TimeoutUpdate),
            Delay:      10 * time.Second,
        }

        _, err = stateConf.WaitForState()
        if err != nil {
            return fmt.Errorf(
                "Error waiting for twcc_firewall %s to become ACTIVE: %v", firewallID, err)
        }
    }

    return resourceFirewallRead(d, meta)
}

func resourceFirewallDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    platform := d.Get("platform").(string)
    firewallID := d.Id()
    resourcePath := fmt.Sprintf("api/v3/%s/firewalls/%s/", platform, firewallID)
    _, err := config.doNormalRequest(platform, resourcePath, "DELETE", nil)

    if err != nil {
        return fmt.Errorf("Unable to delete firewall %s: on %s %v", firewallID, platform, err)
    }

    stateConf := &resource.StateChangeConf{
        Pending:    []string{"DELETING", "PENDING_UPDATE", "PENDING_DELETE"},
        Target:     []string{"DELETED", "ERROR"},
        Refresh:    firewallStateRefreshForDeletedFunc(config, platform, resourcePath),
        Timeout:    d.Timeout(schema.TimeoutDelete),
        Delay:      10 * time.Second,
        MinTimeout: 3 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf( 
            "Error waiting for twcc_firewall %s to become DELETED: %v", firewallID, err)
    }

    d.SetId("")

    return nil
}
