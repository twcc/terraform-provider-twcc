package twcc

import (
    "fmt"
    "encoding/json"
    "strings"

    "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func flattenContainerPortsInfo(v []interface{}) []interface{} {
    portsInfo := make([]interface{}, len(v))
    for i, data := range v {
        port := make(map[string]interface{})
        info := data.(map[string]interface{})
        port["name"] = info["name"].(string)
        port["port"] = int(info["port"].(float64))
        port["protocol"] = info["protocol"].(string)
        portsInfo[i] = port
    }
    return portsInfo
}

func flattenContainerVolumesInfo(v []interface{}) []interface{} {
    volumesInfo := make([]interface{}, len(v))
    for i, data := range v {
        volume := make(map[string]interface{})
        info := data.(map[string]interface{})
        volume["mount_path"] = info["mountPath"].(string)
        volume["path"] = info["path"].(string)
        volume["read_only"] = info["readOnly"].(bool)
        volume["type"] = info["type"].(string)
        volumesInfo[i] = volume
    }
    return volumesInfo
}

func flattenPodContainerInfo(v []interface{}) []interface{} {
    containerInfo := make([]interface{}, len(v))
    for i, data := range v {
        container := make(map[string]interface{})
        info := data.(map[string]interface{})
        imageStringArray := strings.Split(info["image"].(string), "/")
        container["image"] = imageStringArray[len(imageStringArray) - 1]
        container["name"] = info["name"].(string)
        container["ports"] = flattenContainerPortsInfo(info["ports"].([]interface{}))
        container["volumes"] = flattenContainerVolumesInfo(info["volumes"].([]interface{}))
        containerInfo[i] = container
    }
    return containerInfo
}

func flattenServicePortsInfo(v []interface{}) []interface{} {
    portsInfo := make([]interface{}, len(v))
    for i, data := range v {
        port := make(map[string]interface{})
        info := data.(map[string]interface{})
        port["port"] = int(info["port"].(float64))
        port["protocol"] = info["protocol"].(string)
        port["target_port"] = int(info["target_port"].(float64))
        portsInfo[i] = port
    }
    return portsInfo
}

func flattenSitePodInfo(v []interface{}) []interface{} {
    podInfo := make([]interface{}, len(v))
    for i, data := range v {
        pod := make(map[string]interface{})
        info := data.(map[string]interface{})
        pod["container"] = flattenPodContainerInfo(info["container"].([]interface{}))
        pod["flavor"] = info["flavor"].(string)
        pod["message"] = info["message"].(string)
        pod["name"] = info["name"].(string)
        pod["reason"] = info["reason"].(string)
        pod["status"] = info["status"].(string)
        podInfo[i] = pod
    }
    return podInfo
}

func flattenSiteServersInfo(v []interface{}) []interface{} {
    serversInfo := make([]interface{}, len(v))
    for i, data := range v {
        server := make(map[string]interface{})
        info := data.(map[string]interface{})
        server["flavor_id"] = fmt.Sprintf("%d", int(info["flavor_id"].(float64)))
        server["hostname"] = info["hostname"].(string)
        server["id"] = fmt.Sprintf("%d", int(info["id"].(float64)))
        server["status"] = info["status"].(string)
        serversInfo[i] = server
    }
    return serversInfo
}

func flattenSiteServiceInfo(v []interface{}) []interface{} {
    serviceInfo := make([]interface{}, len(v))
    for i, data := range v {
        service := make(map[string]interface{})
        info := data.(map[string]interface{})
        service["name"] = info["name"].(string)
        service["net_type"] = info["net_type"].(string)
        service["ports"] = flattenServicePortsInfo(info["ports"].([]interface{}))
        service["public_ip"] = info["public_ip"].([]interface{})
        serviceInfo[i] = service
    }
    return serviceInfo
}

func siteStateRefreshFunc(
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

func siteStateRefreshForDeletedFunc(
        config *PConfig,
        host string, 
        resourcePath string) resource.StateRefreshFunc {
    return func() (interface{}, string, error) {
        response, err := config.doNormalRequest(host, resourcePath, "GET", nil)

        if err != nil {
            if _, ok := err.(ErrDefault404); ok {
                return response, "Deleted", nil
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
