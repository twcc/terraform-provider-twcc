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

type LoadBalancerCreateBody struct {
    Desc	string		`json:"desc,omitempty"`
    Name	string		`json:"name"`
    Listeners	[]ListenerData	`json:"listeners"`
    Pools	[]PoolData	`json:"pools"`
    PrivateNet	string		`json:"private_net"`
}

type ListenerData struct {
    DefaultTLSContainerRef	string		`json:"default_tls_container_ref,omitempty"`
    Name			string		`json:"name"`
    PoolName			string		`json:"pool_name"`
    Protocol			string		`json:"protocol"`
    ProtocolPort		int		`json:"protocol_port"`
    SNIContainerRefs		[]string	`json:"sni_container_refs,omitempty"`
}

type PoolData struct {
    Delay		int		`json:"delay,omitempty"`
    ExpectedCodes	string		`json:"expected_codes,omitempty"`
    HTTPMethod		string		`json:"http_method,omitempty"`
    MaxRetries          int     	`json:"max_retries,omitempty"`
    Members		*[]MemberData	`json:"members,omitempty"`
    Method		string		`json:"method"`
    Name		string		`json:"name"`
    Protocol		string		`json:"protocol"`
    Timeout		int		`json:"timeout,omitempty"`
    MonitorType		string		`json:"monitor_type,omitempty"`
    URLPath		string		`json:"url_path,omitempty"`
}

type LoadBalancerUpdateBody struct {
    Listeners   []ListenerData  `json:"listeners,omitempty"`
    Pools       []PoolData      `json:"pools,omitempty"`
}

type MemberData struct {
    IP		string	`json:"ip,omitempty"`
    Port	int	`json:"port,omitempty"`
    Weight	int	`json:"weight,omitempty"`
}

func resourceLoadBalancer() *schema.Resource {
    return &schema.Resource{
        Create: resourceLoadBalancerCreate,
        Read:   resourceLoadBalancerRead,
        Update:	resourceLoadBalancerUpdate,
        Delete: resourceLoadBalancerDelete,

        Timeouts: &schema.ResourceTimeout{
            Create: schema.DefaultTimeout(15 * time.Minute),
            Update: schema.DefaultTimeout(15 * time.Minute),
            Delete: schema.DefaultTimeout(15 * time.Minute),
        },

        Schema: map[string]*schema.Schema{
            "active_connections": {
                Type:		schema.TypeInt,
                Computed:	true,
            },

            "create_time": {
                Type:		schema.TypeString,
                Computed:	true,
                ForceNew:	true,
            },

            "desc": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "listeners": {
                Type:		schema.TypeList,
                Required:	true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "default_tls_container_ref": {
                            Type:	schema.TypeString,
                            Optional:	true,
                        },

                        "name": {
                            Type:	schema.TypeString,
                            Required:	true,
                            ForceNew:	true,
                        },

                        "pool": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },

                        "pool_name": {
                            Type:	schema.TypeString,
                            Required:	true,
                        },

                        "protocol": {
                            Type:	schema.TypeString,
                            Required:	true,
                            ForceNew:	true,
                        },

                        "protocol_port": {
                            Type:	schema.TypeInt,
                            Required:	true,
                            ForceNew:	true,
                        },

                        "sni_container_refs": {
                            Type:	schema.TypeList,
                            Optional:	true,
                            Elem:	&schema.Schema{
                                Type:	schema.TypeString,
                            },
                        },

                        "status": {
                            Type:	schema.TypeString,
                            Computed:	true,	
                        },
                    },
                },
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

            "pools": {
                Type:           schema.TypeList,
                Required:       true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "id": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },

                        "members": {
                            Type:	schema.TypeList,
                            Optional:	true,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                     "ip": {
                                         Type:		schema.TypeString,
                                         Required:	true,
                                     },

                                     "port": {       
                                         Type:		schema.TypeInt,
                                         Optional:	true,
                                         Default:	80,
                                     },

                                     "status": {
                                         Type:		schema.TypeString,
                                         Computed:	true,
                                     },

                                     "weight": {     
                                         Type:		schema.TypeInt,
                                         Optional:	true,
                                         Default:	1,
                                     },
                                },
                            },

                            DiffSuppressFunc:	lbMembersDiffFunc,
                        },

                        "method": {
                            Type:	schema.TypeString,
                            Required:	true,
                        },

                        "monitor": {
                            Type:	schema.TypeList,
                            Optional:	true,
                            Computed:	true,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                     "delay": {
                                         Type:		schema.TypeInt,
                                         Optional:	true,
                                         Computed:	true,
                                     },

                                     "expected_codes": {
                                         Type:		schema.TypeString,
                                         Optional:	true,
                                     },

                                     "http_method": {
                                         Type:		schema.TypeString,
                                         Optional:	true,
                                     },

                                     "max_retries": {
                                         Type:		schema.TypeInt,
                                         Optional:	true,
                                         Computed:	true,
                                     },

                                     "monitor_type": {
                                         Type:		schema.TypeString,
                                         Optional:	true,
                                         Computed:	true,
                                     },

                                     "timeout": {
                                         Type:		schema.TypeInt,
                                         Optional:	true,
                                         Computed:	true,
                                     },

                                     "url_path": {
                                         Type:		schema.TypeString,
                                         Optional:	true,
                                     },
                                },
                            },

                            MaxItems:	1,
                        },

                        "name": {
                            Type:	schema.TypeString,
                            Required:	true,
                            ForceNew:	true,
                        },

                        "protocol": {
                            Type:	schema.TypeString,
                            Required:	true,
                            ForceNew:	true,
                        },
 
                        "status": {
                            Type:	schema.TypeString,
                            Computed:	true,
                        },
                    },
                },
            },

            "private_net": {     
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

            "total_connections": {
                Type:		schema.TypeInt,
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

            "vip": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "waf": {
                Type:		schema.TypeMap,
                Computed:	true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
        },
    }
}

func resourceLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    desc := d.Get("desc").(string)  
    listeners := d.Get("listeners").([]interface{})
    listenerArray := make([]ListenerData, len(listeners))
    for i, listener := range listeners {
        l_obj := listener.(map[string]interface{})
        refs := l_obj["sni_container_refs"].([]interface{})
        refArray := make([]string, len(refs))
        for j, ref := range refs{
            refArray[j] = ref.(string)
        }

        listenerBody := ListenerData{
            DefaultTLSContainerRef:	l_obj["default_tls_container_ref"].(string),
            Name:			l_obj["name"].(string),
            PoolName:			l_obj["pool_name"].(string),
            Protocol:			l_obj["protocol"].(string),
            ProtocolPort:		l_obj["protocol_port"].(int),
            SNIContainerRefs:		refArray,
        }

        listenerArray[i] = listenerBody
    }

    name := d.Get("name").(string)
    platform := d.Get("platform").(string)
    pools := d.Get("pools").([]interface{})
    needUpdate := false
    poolArray := make([]PoolData, len(pools))
    for i, pool := range pools {
        p_obj := pool.(map[string]interface{})
        members := p_obj["members"].([]interface{})
        memberArray := make([]MemberData, len(members))
        for j, member := range members{
            member_obj := member.(map[string]interface{})
            memberBody := MemberData {
                IP:	member_obj["ip"].(string),
                Port:	member_obj["port"].(int),
                Weight:	member_obj["weight"].(int),
            }

            memberArray[j] = memberBody
        }

        if len(members) > 0 {
            needUpdate = true
        }

        poolBody := PoolData {
            Members:	&memberArray,
            Method:	p_obj["method"].(string),
            Name:	p_obj["name"].(string),
            Protocol:	p_obj["protocol"].(string),
        }

        monitorArray := p_obj["monitor"].([]interface{})
        if len(monitorArray) != 0 {
            info := monitorArray[0].(map[string]interface{})
            poolBody.Delay = info["delay"].(int)
            poolBody.ExpectedCodes = info["expected_codes"].(string)
            poolBody.HTTPMethod = info["http_method"].(string)
            poolBody.MaxRetries = info["max_retries"].(int)
            poolBody.MonitorType = info["monitor_type"].(string)
            poolBody.Timeout = info["timeout"].(int)
            poolBody.URLPath = info["url_path"].(string)
        }

        poolArray[i] = poolBody
    }

    privateNet := d.Get("private_net").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/loadbalancers/", platform)

    body := LoadBalancerCreateBody {
        Desc:		desc,
        Listeners:	listenerArray,
        Name:		name,
        Pools:		poolArray,
        PrivateNet:	privateNet,
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    response, err := config.doNormalRequest(platform, resourcePath, "POST", buf)

    if err != nil {
        return fmt.Errorf("Error creating twcc_loadbalancer %s on %s: %v", name, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    lbID := int(data["id"].(float64))
    d.SetId(fmt.Sprintf("%d", lbID))

    newPath := fmt.Sprintf("%s/%d/", resourcePath, lbID)
    stateConf := &resource.StateChangeConf{
        Pending:    []string{"BUILD",},
        Target:     []string{"ACTIVE", "DOWN", "ERROR"},
        Refresh:    lbStateRefreshFunc(config, platform, newPath),
        Timeout:    d.Timeout(schema.TimeoutCreate),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf(
            "Error waiting for twcc_loadbalancer %s to become ACTIVE: %v", lbID, err)
    }

    // Update LB if user define member data
    if needUpdate == true {
        body := LoadBalancerUpdateBody {
            Pools:	poolArray,
        }

        buf = new(bytes.Buffer)
        json.NewEncoder(buf).Encode(body)
        _, err = config.doNormalRequest(platform, newPath, "PATCH", buf)

       if err != nil {
            return fmt.Errorf("Error updating twcc_loadbalancer %d on %s: %v", lbID, platform, err)
        }

        stateConf := &resource.StateChangeConf{
            Pending:    []string{"UPDATING"},
            Target:     []string{"ACTIVE", "ERROR"},
            Refresh:    lbStateRefreshFunc(config, platform, newPath),
            Timeout:    d.Timeout(schema.TimeoutUpdate),
            Delay:      10 * time.Second,
        }

        _, err = stateConf.WaitForState()
        if err != nil {
            return fmt.Errorf(
                "Error waiting for twcc_loadbalancer %d to become ACTIVE: %v", lbID, err)
        }
    }

    d.Set("desc", desc)
    d.Set("name", name)
    d.Set("platform", platform)
    d.Set("private_net", privateNet)
    return resourceLoadBalancerRead(d, meta)
}

func resourceLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    lbID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/loadbalancers/%s/", platform, lbID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to retrieve loadbalancer %s on %s: %v", lbID, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve loadbalancer json data: %v", err)
    }

    log.Printf("[DEBUG] Retrieved twcc_loadbalancer %s", d.Id())
    d.Set("active_connections", data["active_connections"])
    d.Set("create_time", data["create_time"])
    listenerInfo := flattenLBListenerInfo(data["listeners"].([]interface{}), data["pools"].([]interface{}))
    d.Set("listeners", listenerInfo)
    poolInfo := flattenLBPoolInfo(data["pools"].([]interface{}))
    d.Set("pools", poolInfo)
    d.Set("status", data["status"])
    d.Set("status_reason", data["status_reason"])
    d.Set("total_connections", data["total_connections"])
    d.Set("user", data["user"].(map[string]interface{}))
    d.Set("vip", data["vip"])
    if waf, ok := data["waf"].(map[string]interface{}); ok {
        waf["id"] = fmt.Sprintf("%d", int(waf["id"].(float64)))
        d.Set("waf", waf)
    } else {
        d.Set("waf", data["waf"])
    }

    return nil
}

func resourceLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    var body LoadBalancerUpdateBody
    if d.HasChange("listeners") || d.HasChange("pools") {
        if d.HasChange("listeners") {
            _, newListeners := d.GetChange("listeners")
            listeners := newListeners.([]interface{})
            listenerArray := make([]ListenerData, len(listeners))
            for i, listener := range listeners {
                l_obj := listener.(map[string]interface{})
                refs := l_obj["sni_container_refs"].([]interface{})
                refArray := make([]string, len(refs))
                for j, ref := range refs{
                    refArray[j] = ref.(string)
                }

                listenerBody := ListenerData{
                    DefaultTLSContainerRef:	l_obj["default_tls_container_ref"].(string),
                    Name:			l_obj["name"].(string),
                    PoolName:			l_obj["pool_name"].(string),
                    Protocol:			l_obj["protocol"].(string),
                    ProtocolPort:		l_obj["protocol_port"].(int),
                    SNIContainerRefs:		refArray,
                }

                listenerArray[i] = listenerBody
            }

            body.Listeners = listenerArray
        }

        if d.HasChange("pools") {
            _, newPools := d.GetChange("pools")
            pools := newPools.([]interface{})
            poolArray := make([]PoolData, len(pools))
            for i, pool := range pools {
                p_obj := pool.(map[string]interface{})
                members := p_obj["members"].([]interface{})
                memberArray := make([]MemberData, len(members))
                for j, member := range members{
                    member_obj := member.(map[string]interface{})
                    memberBody := MemberData {
                        IP:	member_obj["ip"].(string),
                        Port:	member_obj["port"].(int),
                        Weight:	member_obj["weight"].(int),
                    }

                    memberArray[j] = memberBody
                }

                poolBody := PoolData {
                    Members:	&memberArray,
                    Method:		p_obj["method"].(string),
                    Name:		p_obj["name"].(string),
                    Protocol:	p_obj["protocol"].(string),
                }

                poolArray[i] = poolBody
            }

            body.Pools = poolArray
        }

        lbID := d.Id()
        platform := d.Get("platform").(string)
        resourcePath := fmt.Sprintf("api/v3/%s/loadbalancers/%s/", platform, lbID)

        buf := new(bytes.Buffer)
        json.NewEncoder(buf).Encode(body)
        _, err := config.doNormalRequest(platform, resourcePath, "PATCH", buf)

        if err != nil {
            return fmt.Errorf("Error updating twcc_loadbalancer %s on %s: %v", lbID, platform, err)
        }

        stateConf := &resource.StateChangeConf{
            Pending:	[]string{"UPDATING"},
            Target:	[]string{"ACTIVE", "ERROR"},
            Refresh:	lbStateRefreshFunc(config, platform, resourcePath),
            Timeout:	d.Timeout(schema.TimeoutUpdate),
            Delay:	10 * time.Second,
        }

        _, err = stateConf.WaitForState()
        if err != nil {
            return fmt.Errorf(
                "Error waiting for twcc_loadbalancer %s to become ACTIVE: %v", lbID, err)
        }
    }

    return resourceLoadBalancerRead(d, meta)
}

func resourceLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    platform := d.Get("platform").(string)
    lbID := d.Id()
    resourcePath := fmt.Sprintf("api/v3/%s/loadbalancers/%s/", platform, lbID)
    _, err := config.doNormalRequest(platform, resourcePath, "DELETE", nil)

    if err != nil {
        return fmt.Errorf("Unable to delete loadbalancer %s: on %s %v", lbID, platform, err)
    }

    stateConf := &resource.StateChangeConf{
        Pending:    []string{"DELETING"},
        Target:     []string{"DELETED", "ERROR"},
        Refresh:    lbStateRefreshForDeletedFunc(config, platform, resourcePath),
        Timeout:    d.Timeout(schema.TimeoutDelete),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf( 
            "Error waiting for twcc_loadbalancer %s to become DELETED: %v", lbID, err)
    }

    d.SetId("")

    return nil
}
