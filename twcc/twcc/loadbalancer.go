package twcc

import (
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func checkIsChanged(oldArray []interface{}, newArray []interface{}) bool {
    if len(oldArray) != len(newArray){
       return false
    }

    equalCount := 0
    for _, x := range newArray{
        newObject := x.(map[string]interface{})
        newIP := newObject["ip"].(string)
        newPort := newObject["port"].(int)
        newWeight := newObject["weight"].(int)
        for _, y := range oldArray {
            oldObject := y.(map[string]interface{})
            oldIP := oldObject["ip"].(string)
            oldPort := oldObject["port"].(int)
            oldWeight := oldObject["weight"].(int)
            if newIP == oldIP && newPort == oldPort && newWeight == oldWeight{
                equalCount += 1
                break
            }
        }
    }
    return equalCount == len(oldArray)
}

func lbMembersDiffFunc(k, old, new string, d *schema.ResourceData) bool {
    oldPools, newPools := d.GetChange("pools")
    oldPoolArray := oldPools.([]interface{})
    newPoolArray := newPools.([]interface{})
    if len(oldPoolArray) == 0 && len(newPoolArray) != 0 {
        return false
    }
    for _, newPool := range oldPoolArray {
        newPoolDict := newPool.(map[string]interface{})
        newPoolName := newPoolDict["name"].(string)
        newPoolMembers := newPoolDict["members"].([]interface{})
        for _, oldPool := range newPoolArray{
            oldPoolDict := oldPool.(map[string]interface{})
            oldPoolName := oldPoolDict["name"].(string)
            if oldPoolName == newPoolName {
                oldPoolMembers := oldPoolDict["members"].([]interface{})
                if !(checkIsChanged(oldPoolMembers, newPoolMembers)) {
                    return false
                }
            }
        }
    }
    return true
}

func flattenLBListenerInfo(v []interface{}, v2 []interface{}) []interface{} {
    listenerInfo := make([]interface{}, len(v))
    for i, data := range v {
        listener := make(map[string]interface{})
        info := data.(map[string]interface{})
        if defaultRef, ok := info["default_tls_container_ref"].(float64); ok {
            listener["default_tls_container_ref"] = fmt.Sprintf("%d", int(defaultRef))
        } else {
            listener["default_tls_container_ref"] = ""
        }
        listener["name"] = info["name"].(string)
        pool_id := info["pool"].(float64)
        for _, data2 := range v2 {
            info2 := data2.(map[string]interface{})
            if info2["id"].(float64) == pool_id {
                listener["pool_name"] = info2["name"].(string)
                break
            }
        }
        listener["pool"] = fmt.Sprintf("%d", int(pool_id))
        listener["protocol"] = info["protocol"].(string)
        listener["protocol_port"] = int(info["protocol_port"].(float64))
        refs := info["sni_container_refs"].([]interface{})
        s_refs := make([]string, len(refs))
        for j, ref := range refs {
            s_refs[j] = fmt.Sprintf("%d", int(ref.(float64)))
        }
        listener["sni_container_refs"] = s_refs
        listener["status"] = info["status"].(string)
        listenerInfo[i] = listener
    }
    return listenerInfo
}

func flattenLBMemberInfo(v []interface{}) []interface{} {
    memberInfo := make([]interface{}, len(v))
    for i, data := range v {
        member := make(map[string]interface{})
        info := data.(map[string]interface{})
        member["ip"] = info["ip"].(string)
        member["port"] = int(info["port"].(float64))
        member["weight"] = int(info["weight"].(float64))
        memberInfo[i] = member
    }
    return memberInfo
}

func flattenLBMonitorInfo(v map[string]interface{}) []interface{} {
    monitorInfo := make([]interface{}, 1)
    info := make(map[string]interface{})
    info["delay"] = int(v["delay"].(float64))
    if expectedCodes, ok := v["expected_codes"].(string); ok {
        info["expected_codes"] = expectedCodes
    } else {
        info["expected_codes"] = ""
    }

    if httpMethod, ok := v["http_method"].(string); ok {
        info["http_method"] = httpMethod
    } else {
        info["http_method"] = ""
    }

    info["max_retries"] = int(v["max_retries"].(float64))
    info["monitor_type"] = v["monitor_type"].(string)
    info["timeout"] = int(v["timeout"].(float64))
    if urlPath, ok := v["url_path"].(string); ok {
        info["url_path"] = urlPath
    } else {
        info["url_path"] = ""
    }

    monitorInfo[0] = info
    return monitorInfo
}

func flattenLBPoolInfo(v []interface{}) []interface{} {
    poolInfo := make([]interface{}, len(v))
    for i, data := range v {
        pool := make(map[string]interface{})
        info := data.(map[string]interface{})
        pool["id"] = fmt.Sprintf("%d", int(info["id"].(float64)))
        members := info["members"].([]interface{})
        pool["members"] = flattenLBMemberInfo(members)
        pool["method"] = info["method"].(string)
        if monitor, ok := info["monitor"].(map[string]interface{}); ok {
            monitorInfo := flattenLBMonitorInfo(monitor)
            pool["monitor"] = monitorInfo
        } else {
            pool["monitor"] = make([]interface{}, 0)
        }
        pool["name"] = info["name"].(string)
        pool["protocol"] = info["protocol"].(string)
        pool["status"] = info["status"].(string)
        poolInfo[i] = pool
    }   
    return poolInfo
}

func lbStateRefreshFunc(
        config *PConfig,
        host string,
        resourcePath string) resource.StateRefreshFunc {
    return func() (interface{}, string, error) {
        response, err := config.doNormalRequest(host, resourcePath, "GET", nil)
        if err != nil {
            return nil, "", err
        }

        var data map[string]interface{}
        err = json.Unmarshal([]byte(response), &data)

        if err != nil {
            return nil, "", err
        }

        return data, data["status"].(string), nil
    }
}

func lbStateRefreshForDeletedFunc(
        config *PConfig,
        host string, 
        resourcePath string) resource.StateRefreshFunc {
    return func() (interface{}, string, error) {
        response, err := config.doNormalRequest(host, resourcePath, "GET", nil)

        if err != nil {
            if _, ok := err.(ErrDefault404); ok {
                return response, "DELETED", nil
            }
            return response, "", err
        }

        var data map[string]interface{}
        err = json.Unmarshal([]byte(response), &data)

        if err != nil {
            return nil, "", err
        }

        return data, data["status"].(string), nil
    }
}
