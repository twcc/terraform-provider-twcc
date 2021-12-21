package twcc

import (
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func flattenVolumeHostInfo(v interface{}) map[string]string {
    host := make(map[string]string)
    info := v.(map[string]interface{})
    host["hostname"] = info["hostname"].(string)
    host["id"] = fmt.Sprintf("%d", int(info["id"].(float64)))
    return host
}

func flattenVolumeSnapshotsInfo(v []interface{}) []interface{} {
    snapshotsInfo := make([]interface{}, len(v))
    for i, obj := range v {
        snapshot := fmt.Sprintf("%d", int(obj.(float64)))
        snapshotsInfo[i] = snapshot
    }
    return snapshotsInfo
}

func volumeStateRefreshFunc(
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

func volumeStateRefreshForDeletedFunc(
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
