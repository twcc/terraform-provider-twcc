package twcc

import (
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)


func flattenVPNConnectionInfo(connection interface{}) (VPNConnectionInfo []map[string]interface{}) {
    if vpn_connection, ok := connection.(map[string]interface{}); ok {
        VPNConnectionInfo = append(VPNConnectionInfo, vpn_connection)
    }
    return VPNConnectionInfo
}


func connectionStateRefreshFunc(
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

        connectionData := data["connection"]
        if connectionData != nil {
            return connectionData, "CREATED", nil
        } else {
            return data, "CREATING", nil
        }
    }
}

func connectionStateRefreshForDeletedFunc(
        config *PConfig,
        host string, 
        resourcePath string) resource.StateRefreshFunc {
    return func() (interface{}, string, error) {
        response, err := config.doNormalRequest(host, resourcePath, "GET", nil)

        if err != nil {
            return response, "", err
        }

        var data map[string]interface{}
        err = json.Unmarshal([]byte(response), &data)

        if err != nil {
            return nil, "", err
        }

        connectionData := data["connection"]
        if connectionData == nil {
            return data, "DELETED", nil
        } else {
            return connectionData, "DELETING", nil
        }
    }
}
