package twcc

import (
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func flattenFirewallObjectInfo(v []interface{}) []interface{} {
    networkInfo := make([]interface{}, len(v))
    for i, obj := range v {
        info := obj.(map[string]interface{})
        networkInfo[i] = fmt.Sprintf("%d", int(info["id"].(float64)))
    }
    return networkInfo
}

func firewallStateRefreshFunc(
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

func firewallStateRefreshForDeletedFunc(
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
