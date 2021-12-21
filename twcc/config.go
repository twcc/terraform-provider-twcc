package twcc

import (
    "bytes"
    "fmt"
)

type Config struct {
    TWCC_APIKEY		string
    APIGW_URL		string

    APIGWClient		*ProviderClient
}

func (c *Config) LoadAndValidate() error {
    if c.TWCC_APIKEY == "" {
        return fmt.Errorf("'TWCC_APIKEY' must be specified")
    }

    if c.APIGW_URL == "" {
        return fmt.Errorf("'APIGW_URL' must be specified")
    }

    client := newProviderClient(c.TWCC_APIKEY, c.APIGW_URL)
    c.APIGWClient = client

    return nil
}

func (c *Config) doNormalRequest (
        resourceHost string,
        resourcePath string,
        method string,
        body *bytes.Buffer) (string, error) {
    response, err := c.APIGWClient.doRequest(resourceHost, resourcePath, method, body, nil)
    return response, err
}

func (c *Config) doCreateSiteRequest (
        resourceHost string,
        resourcePath string,
        method string,
        body *bytes.Buffer,
        headers map[string]string) (string, error) {
    response, err := c.APIGWClient.doRequest(resourceHost, resourcePath, method, body, headers)
    return response, err
}
