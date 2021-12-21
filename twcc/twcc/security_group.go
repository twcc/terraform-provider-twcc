package twcc

import (
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


func flattenSecurityGroupRulesInfo(v []interface{}) []interface{} {
    security_group_rules := make([]interface{}, len(v))
    for i, data := range v {
        security_group_rule := make(map[string]interface{})
        info := data.(map[string]interface{})
        security_group_rule["id"] = info["id"].(string)
        security_group_rule["direction"] = info["direction"].(string)
        security_group_rule["ethertype"] = info["ethertype"].(string)
        security_group_rule["remote_ip_prefix"] = info["remote_ip_prefix"].(string)
        security_group_rule["protocol"] = info["protocol"].(string)
        if port_range_min, ok := info["port_range_min"].(float64); ok {
            security_group_rule["port_range_min"] = int(port_range_min)
        }
        if port_range_max, ok := info["port_range_max"].(float64); ok {
            security_group_rule["port_range_max"] = int(port_range_max)
        }
        security_group_rules[i] = security_group_rule
    }
    return security_group_rules
}


func foundSecurityGroupRule(v map[string]interface{}, d *schema.ResourceData) bool {
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
    return !(
        (direction != "" && v["direction"] != direction) ||
        (direction == "" && v["direction"] != "ingress") ||
        (protocol != "" && v["protocol"] != protocol) ||
        (protocol == "" && v["protocol"] != "tcp") ||
        (remote_ip_prefix != "" && v["remote_ip_prefix"] != remote_ip_prefix) ||
        (remote_ip_prefix == "" && v["remote_ip_prefix"] != "0.0.0.0/0") ||
        (port_range_min != 0 && int(v["port_range_min"].(float64)) != port_range_min) ||
        (port_range_min == 0 && v["port_range_min"] != nil) ||
        (port_range_max != 0 && int(v["port_range_max"].(float64)) != port_range_max) ||
        (port_range_max == 0 && v["port_range_max"] != nil))
}
