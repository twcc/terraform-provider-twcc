package twcc

import (
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func relationStateRefreshFunc(
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

        relationData := data["auto_scaling_policy"]
        if relationData != nil {
            relationDict := relationData.(map[string]interface{})
            status := relationDict["status"].(string)
            return relationDict, status, nil
        } else {
            return data, "ASSOCIATING", nil
        }
    }
}

func relationStateRefreshForDeletedFunc(
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

        relationData := data["auto_scaling_policy"]
        if relationData == nil {
            return data, "DELETED", nil
        } else {
            return relationData, "DISASSOCIATING", nil
        }
    }
}
