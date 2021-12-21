package twcc

import (
    "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func firewallRuleStateRefreshForDeletedFunc(
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
        } else {
            return response, "DELETING", nil
        }
    }
}
