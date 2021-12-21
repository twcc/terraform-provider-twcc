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

type VPNCreateBody struct {
    IKEPolicy		string	`json:"ike_policy"`
    IPSecPolicy		string	`json:"ipsec_policy"`
    Name		string	`json:"name"`
    PrivateNetwork	string	`json:"private_network"`
}

func resourceVPN() *schema.Resource {
    return &schema.Resource{
        Create: resourceVPNCreate,
        Read:   resourceVPNRead,
        Delete: resourceVPNDelete,

        Timeouts: &schema.ResourceTimeout{
            Delete: schema.DefaultTimeout(10 * time.Minute),
        },

        Schema: map[string]*schema.Schema{
            "ike_policy": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "ipsec_policy": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "local_address": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "local_cidr": {
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

            "private_network": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "status": {
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

func resourceVPNCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    ikePolicy := d.Get("ike_policy").(string)
    ipsecPolicy := d.Get("ipsec_policy").(string)
    name := d.Get("name").(string)
    platform := d.Get("platform").(string)
    privateNetwork := d.Get("private_network").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/vpn_services/", platform)

    body := VPNCreateBody {
        IKEPolicy:	ikePolicy,
        IPSecPolicy:	ipsecPolicy,
        Name:		name,
        PrivateNetwork:	privateNetwork,
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    response, err := config.doNormalRequest(platform, resourcePath, "POST", buf)

    if err != nil {
        return fmt.Errorf("Error creating twcc_vpn %s on %s: %v", name, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    vpnID := int(data["id"].(float64))
    d.SetId(fmt.Sprintf("%d", vpnID))
    d.Set("ike_policy", ikePolicy)
    d.Set("ipsec_policy", ipsecPolicy)
    d.Set("name", name)
    d.Set("platform", platform)
    d.Set("private_network", privateNetwork)
    return resourceVPNRead(d, meta)
}

func resourceVPNRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    vpnID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/vpn_services/%s/", platform, vpnID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to retrieve vpn %s on %s: %v", vpnID, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve vpn json data: %v", err)
    }

    log.Printf("[DEBUG] Retrieved twcc_vpn %s", d.Id())
    d.Set("local_address", data["local_address"])
    d.Set("local_cidr", data["local_cidr"])
    d.Set("status", data["status"])
    d.Set("user", data["user"])
    return nil
}

func resourceVPNDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    vpnID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/vpn_services/%s/", platform, vpnID)
    _, err := config.doNormalRequest(platform, resourcePath, "DELETE", nil)

    if err != nil {
        return fmt.Errorf("Unable to delete vpn %s: on %s %v", vpnID, platform, err)
    }

    stateConf := &resource.StateChangeConf{
        Pending:    []string{"DELETING"},
        Target:     []string{"DELETED"},
        Refresh:    vpnStateRefreshForDeletedFunc(config, platform, resourcePath),
        Timeout:    d.Timeout(schema.TimeoutDelete),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf(
            "Error waiting for twcc_vpn %s to become DELETED: %v", vpnID, err)
    }

    d.SetId("")

    return nil
}
