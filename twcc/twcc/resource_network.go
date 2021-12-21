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

type NetworkCreateBody struct {
    CIDR	string		`json:"cidr"`
    DNSDomain	string		`json:"dns_domain,omitempty"`
    Gateway	string		`json:"gateway"`
    Name	string		`json:"name"`
    Nameservers	[]string	`json:"nameservers,omitempty"`
    Project	string		`json:"project"`
    WithRouter	bool		`json:"with_router,omitempty"`
}

func resourceNetwork() *schema.Resource {
    return &schema.Resource{
        Create: resourceNetworkCreate,
        Read:   resourceNetworkRead,
        Delete: resourceNetworkDelete,

        Timeouts: &schema.ResourceTimeout{
            Create: schema.DefaultTimeout(10 * time.Minute),
            Delete: schema.DefaultTimeout(10 * time.Minute),
        },

        Schema: map[string]*schema.Schema{
            "cidr": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "create_time": {
                Type:		schema.TypeString,
                Computed:	true,
                ForceNew:	true,
            },

            "dns_domain": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "ext_net": {
                Type:		schema.TypeString,
                Computed:	true,
                ForceNew:	true,
            },

            "firewall": {
                Type:		schema.TypeMap,
                Computed:	true,
                ForceNew:	true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },

            "gateway": {           
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "ip_version": {        
                Type:		schema.TypeInt,
                Computed:	true,
                ForceNew:	true,
            },

            "name": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:       true,
            },

            "nameservers": {
                Type:		schema.TypeList,
                Optional:	true,
                Computed:	true,
                ForceNew:	true,
                Elem: &schema.Schema{
                    Type:	schema.TypeString,
                },
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

            "with_router": {
                Type:		schema.TypeBool,
                Optional:	true,
                Computed:	true,
                ForceNew:	true,
            },
        },
    }
}

func resourceNetworkCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    cidr := d.Get("cidr").(string)
    dnsDomain := d.Get("dns_domain").(string)
    gateway := d.Get("gateway").(string)
    name := d.Get("name").(string)
    nameServers := d.Get("nameservers").([]interface{})
    platform := d.Get("platform").(string)
    project := d.Get("project").(string)
    withRouter := d.Get("with_router").(bool)
    resourcePath := fmt.Sprintf("api/v3/%s/networks/", platform)

    body := NetworkCreateBody {
        CIDR:		cidr,
        DNSDomain:	dnsDomain,
        Gateway:	gateway,
        Name:		name,
        Nameservers:	flattenNetworkNameServersInfo(nameServers),
        Project:	project,
        WithRouter:	withRouter,
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    response, err := config.doNormalRequest(platform, resourcePath, "POST", buf)

    if err != nil {
        return fmt.Errorf("Error creating twcc_network %s on %s: %v", name, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    networkID := int(data["id"].(float64))
    d.SetId(fmt.Sprintf("%d", networkID))

    newPath := fmt.Sprintf("%s%d/", resourcePath, networkID)
    stateConf := &resource.StateChangeConf{
        Pending:    []string{"BUILD"},
        Target:     []string{"ACTIVE", "ERROR"},
        Refresh:    networkStateRefreshFunc(config, platform, newPath),
        Timeout:    d.Timeout(schema.TimeoutCreate),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf(
            "Error waiting for twcc_network %d to become ACTIVE: %v", networkID, err)
    }

    d.Set("name", name)
    d.Set("platform", platform)
    d.Set("project", project)
    return resourceNetworkRead(d, meta)
}

func resourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    networkID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/networks/%s/", platform, networkID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to retrieve network %s on %s: %v", networkID, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve network json data: %v", err)
    }

    log.Printf("[DEBUG] Retrieved twcc_network %s", d.Id())
    d.Set("cidr", data["cidr"])
    d.Set("create_time", data["create_time"])
    if data["dns_domain"] != nil {
        d.Set("dns_domain", data["dns_domain"])
    }
    d.Set("ext_net", data["ext_net"])
    if firewall, ok := data["firewall"].(map[string]interface{}); ok {
        firewall["id"] = fmt.Sprintf("%v", firewall["id"])
        d.Set("firewall", firewall)
    } else {
        d.Set("firewall", data["firewall"])
    }
    d.Set("ip_version", data["ip_version"])
    d.Set("nameservers", data["nameservers"])
    d.Set("status", data["status"])
    d.Set("status_reason", data["status_reason"])
    d.Set("user", data["user"])
    d.Set("with_router", data["with_router"])
    return nil
}

func resourceNetworkDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    networkID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/networks/%s/", platform, networkID)
    _, err := config.doNormalRequest(platform, resourcePath, "DELETE", nil)

    if err != nil {
        return fmt.Errorf("Unable to delete network %s: on %s %v", networkID, platform, err)
    }

    stateConf := &resource.StateChangeConf{
        Pending:    []string{"DELETING"},
        Target:     []string{"DELETED", "ERROR"},
        Refresh:    networkStateRefreshForDeletedFunc(config, platform, resourcePath),
        Timeout:    d.Timeout(schema.TimeoutDelete),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf( 
            "Error waiting for twcc_network %s to become DELETED: %v", networkID, err)
    }

    d.SetId("")

    return nil
}
