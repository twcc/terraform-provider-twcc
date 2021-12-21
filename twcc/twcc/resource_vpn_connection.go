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

type ConnectionCreateBody struct {
    DPDAction	string		`json:"dpd_action,omitempty"`
    DPDInterval	int		`json:"dpd_interval,omitempty"`
    DPDTimeout	int		`json:"dpd_timeout,omitempty"`
    Initiator	string		`json:"initiator,omitempty"`
    MTU		int		`json:"mtu,omitempty"`
    PeerAddress	string		`json:"peer_address"`
    PeerCIDRs	[]string	`json:"peer_cidrs"`
    PeerID	string		`json:"peer_id,omitempty"`
    PSK		string		`json:"psk,omitempty"`
}

func resourceVPNConnection() *schema.Resource {
    return &schema.Resource{
        Create: resourceVPNConnectionCreate,
        Read:   resourceVPNConnectionRead,
        Delete: resourceVPNConnectionDelete,

        Timeouts: &schema.ResourceTimeout{
            Delete: schema.DefaultTimeout(10 * time.Minute),
        },

        Schema: map[string]*schema.Schema{
            "dpd_action": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
                Default:	"hold",
            },

            "dpd_interval": {
                Type:		schema.TypeInt,
                Optional:	true,
                ForceNew:	true,
                Default:	30,
            },

            "dpd_timeout": {
                Type:		schema.TypeInt,
                Optional:	true,
                ForceNew:	true,
                Default:	120,
            },

            "initiator": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
                Default:	"bi-directional",
            },

            "mtu": {
                Type:		schema.TypeInt,
                Optional:	true,
                ForceNew:       true,
                Default:	1500,
            },

            "peer_address": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "peer_cidrs": {
                Type:		schema.TypeList,
                Required:	true,
                ForceNew:	true,
                Elem: &schema.Schema{
                    Type:	schema.TypeString,
                },
            },

            "peer_id": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "platform": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:       true,
            },

            "psk": {
                Type:           schema.TypeString,
                Required:       true,
                ForceNew:       true,
            },

            "status": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "vpn": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },
        },
    }
}

func resourceVPNConnectionCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    dpdAction := d.Get("dpd_action").(string)
    dpdInterval := d.Get("dpd_interval").(int)
    dpdTimeout := d.Get("dpd_timeout").(int)
    initiator := d.Get("initiator").(string)
    mtu := d.Get("mtu").(int)
    peerAddress := d.Get("peer_address").(string)
    peerCIDRs := d.Get("peer_cidrs").([]interface{})
    peerArray := make([]string, len(peerCIDRs))
    for i, v := range peerCIDRs {
        peerArray[i] = v.(string)
    }

    peerID := d.Get("peer_id").(string)
    platform := d.Get("platform").(string)
    psk := d.Get("psk").(string)
    vpn := d.Get("vpn").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/vpn_services/%s/connection/", platform, vpn)

    body := ConnectionCreateBody {
        DPDAction:	dpdAction,
        DPDInterval:	dpdInterval,
        DPDTimeout:	dpdTimeout,
        Initiator:	initiator,
        MTU:		mtu,
        PeerAddress:	peerAddress,
        PeerCIDRs:	peerArray,
        PeerID:		peerID,
        PSK:		psk,
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    _, err := config.doNormalRequest(platform, resourcePath, "POST", buf)

    if err != nil {
        return fmt.Errorf("Error creating twcc_vpn_connection %s on %s: %v", vpn, platform, err)
    }

    d.SetId(fmt.Sprintf("%s-connection", vpn))

    newPath := fmt.Sprintf("api/v3/%s/vpn_services/%s/", platform, vpn)
    stateConf := &resource.StateChangeConf{
        Pending:    []string{"CREATING"},
        Target:     []string{"CREATED"},
        Refresh:    connectionStateRefreshFunc(config, platform, newPath),
        Timeout:    d.Timeout(schema.TimeoutCreate),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf(
            "Error waiting for twcc_vpn_connection %s to become Ready: %v", vpn, err)
    }

    d.Set("peer_address", peerAddress)
    d.Set("platform", platform)
    d.Set("psk", psk)
    return resourceVPNConnectionRead(d, meta)
}

func resourceVPNConnectionRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    vpnID := d.Get("vpn").(string)
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

    log.Printf("[DEBUG] Retrieved twcc_vpn_connection %s", d.Id())
    if connection, ok := data["connection"].(map[string]interface{}); ok {
        d.Set("dpd_action", connection["dpd_action"])
        d.Set("dpd_interval", connection["dpd_interval"])
        d.Set("dpd_timeout", connection["dpd_timeout"])
        d.Set("initiator", connection["initiator"])
        d.Set("mtu", connection["mtu"])
        d.Set("peer_id", connection["peer_id"])
        d.Set("peer_cidrs", connection["peer_cidrs"])
        d.Set("status", connection["status"])
    } else {
        return fmt.Errorf("VPN connection not found")
    }
    return nil
}

func resourceVPNConnectionDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    vpnID := d.Get("vpn").(string)
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/vpn_services/%s/connection/", platform, vpnID)
    _, err := config.doNormalRequest(platform, resourcePath, "DELETE", nil)

    if err != nil {
        return fmt.Errorf("Unable to delete vpn connection %s: on %s %v", vpnID, platform, err)
    }

    newPath := fmt.Sprintf("api/v3/%s/vpn_services/%s/", platform, vpnID)
    stateConf := &resource.StateChangeConf{
        Pending:    []string{"DELETING"},
        Target:     []string{"DELETED"},
        Refresh:    connectionStateRefreshForDeletedFunc(config, platform, newPath),
        Timeout:    d.Timeout(schema.TimeoutDelete),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf(
            "Error waiting for twcc_vpn_connection %s to become DELETED: %v", vpnID, err)
    }

    d.SetId("")

    return nil
}
