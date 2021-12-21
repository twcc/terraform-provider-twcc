package twcc

import (
    "bytes"
    "crypto/tls"
    "io/ioutil"
    "net/http"
    "strings"
)

type ProviderClient struct {
    HTTPClient	http.Client
    Key		string
    Url 	string
}

func newProviderClient(key string, url string) *ProviderClient {
    pc := ProviderClient{}
    pc.Key = key
    pc.Url = url
    is_https := strings.Contains(url, "https")

    if is_https {
        tr := &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
        client := http.Client{Transport: tr}
        pc.HTTPClient = client
    } else{
        client := http.Client{}
        pc.HTTPClient = client
    }

    return &pc
}

func (pc *ProviderClient) doRequest(
        resourceHost string,
        resourcePath string,
        method string,
        body *bytes.Buffer,
        headers map[string]string) (string, error) {
    url :=  pc.Url + resourcePath
    var req *http.Request
    var err error

    if body != nil {
        req, err = http.NewRequest(method, url, body)
        req.Header.Set("Content-Type", "application/json")
    } else {
        req, err = http.NewRequest(method, url, nil)
    }

    if headers != nil {
        for key, value := range headers {
            req.Header.Set(key, value)
        }
    }

    req.Header.Set("x-api-host", resourceHost)
    req.Header.Set("x-api-key", pc.Key)
    if err != nil {
        return "Initial HTTP request failed", err
    }

    resp, err := pc.HTTPClient.Do(req)
    if err != nil {
        return "Sending request failed", err
    }

    okc := defaultOkCodes(method)
    var ok bool
    for _, code := range okc {
        if resp.StatusCode == code {
            ok = true
            break
        }
    }

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()

    if err != nil {
        return "Read response body failed", err
    }

    if !ok {
        respErr := ErrUnexpectedResponseCode{
            URL:            url,
            Method:         method,
            Expected:       okc,
            Actual:         resp.StatusCode,
            Body:           bodyBytes,
            ResponseHeader: resp.Header,
        }
        switch resp.StatusCode {
        case http.StatusBadRequest:
            err = ErrDefault400{respErr}
            if error400er, ok := err.(Err400er); ok {
                err = error400er.Error400(respErr)
            }
        case http.StatusUnauthorized:
            err = ErrDefault401{respErr}
            if error401er, ok := err.(Err401er); ok {
                err = error401er.Error401(respErr)
            }
        case http.StatusForbidden:
            err = ErrDefault403{respErr}
            if error403er, ok := err.(Err403er); ok {
                err = error403er.Error403(respErr)
            }
        case http.StatusNotFound:
            err = ErrDefault404{respErr}
            if error404er, ok := err.(Err404er); ok {
                err = error404er.Error404(respErr)
            }
        case http.StatusConflict:
            err = ErrDefault409{respErr}
            if error409er, ok := err.(Err409er); ok {
                err = error409er.Error409(respErr)
            }
        case http.StatusInternalServerError:
            err = ErrDefault500{respErr}
            if error500er, ok := err.(Err500er); ok {
                err = error500er.Error500(respErr)
            }
        case http.StatusServiceUnavailable:
            err = ErrDefault503{respErr}
            if error503er, ok := err.(Err503er); ok {
                err = error503er.Error503(respErr)
            }
        }

        if err == nil {
            err = respErr
        }
    }

    return string(bodyBytes), err
}

func defaultOkCodes(method string) []int {
    switch method {
    case "GET", "HEAD":
        return []int{200}
    case "POST":
        return []int{201, 202}
    case "PUT":
        return []int{200, 201, 202}
    case "PATCH":
        return []int{200, 201, 202, 204}
    case "DELETE":
        return []int{202, 204}
    }

    return []int{}
}
