package twcc

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type IKEPolicyCreateBody struct {
    AuthAlgorithm	string		`json:"auth_algorithm,omitempty"`
    EncryptionAlgorithm	string		`json:"encryption_algorithm,omitempty"`
    IKEVersion		string		`json:"ike_version,omitempty"`
    Lifetime		int		`json:"lifetime,omitempty"`
    Name		string		`json:"name"`
    PFS			string		`json:"pfs,omitempty"`
    Project		string		`json:"project"`
}

func resourceIKEPolicy() *schema.Resource {
    return &schema.Resource{
        Create: resourceIKEPolicyCreate,
        Read:   resourceIKEPolicyRead,
        Delete: resourceIKEPolicyDelete,

        Timeouts: &schema.ResourceTimeout{
            Delete: schema.DefaultTimeout(3 * time.Minute),
        },

        Schema: map[string]*schema.Schema{
            "auth_algorithm": {
                Type:		schema.TypeString,
                Optional:	true,
                Computed:	true,
                ForceNew:	true,
            },

            "encryption_algorithm": {
                Type:		schema.TypeString,
                Optional:	true,
                Computed:	true,
                ForceNew:	true,
            },

            "ike_version": {
                Type:		schema.TypeString,
                Optional:	true,
                Computed:	true,
                ForceNew:	true,
            },

            "lifetime": {
                Type:		schema.TypeInt,
                Optional: 	true,
                Computed:	true,
                ForceNew:	true,
            },

            "name": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:       true,
            },

            "pfs": {
                Type:		schema.TypeString,
                Optional:	true,
                Computed:	true,
                ForceNew:	true,
            },

            "platform": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:       true,
            },

            "project": {     
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "user": {       
                Type:		schema.TypeMap,
                Computed:	true,
                ForceNew:	true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
        },
    }
}

func resourceIKEPolicyCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    authAlgorithm := d.Get("auth_algorithm").(string)
    encryptionAlgorithm := d.Get("encryption_algorithm").(string)
    ikeVersion := d.Get("ike_version").(string)
    lifetime := d.Get("lifetime").(int)
    name := d.Get("name").(string)
    pfs := d.Get("pfs").(string)
    platform := d.Get("platform").(string)
    project := d.Get("project").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/ike_policies/", platform)

    body := IKEPolicyCreateBody {
        AuthAlgorithm:		authAlgorithm,
        EncryptionAlgorithm:	encryptionAlgorithm,
        IKEVersion:		ikeVersion,
        Lifetime:		lifetime,
        Name:			name,
        PFS:			pfs,
        Project:		project,
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    response, err := config.doNormalRequest(platform, resourcePath, "POST", buf)

    if err != nil {
        return fmt.Errorf("Error creating twcc_ike_policy %s on %s: %v", name, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    policyID := int(data["id"].(float64))
    d.SetId(fmt.Sprintf("%d", policyID))
    d.Set("name", name)
    d.Set("platform", platform)
    d.Set("project", project)
    return resourceIKEPolicyRead(d, meta)
}

func resourceIKEPolicyRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    policyID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/ike_policies/%s/", platform, policyID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to retrieve IKE policy %s on %s: %v", policyID, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve IKE policy json data: %v", err)
    }

    log.Printf("[DEBUG] Retrieved twcc_ike_policy %s", d.Id())
    d.Set("auth_algorithm", data["auth_algorithm"])
    d.Set("encryption_algorithm", data["encryption_algorithm"])
    d.Set("ike_version", data["ike_version"])
    d.Set("lifetime", data["lifetime"])
    d.Set("pfs", data["pfs"])
    d.Set("user", data["user"])
    return nil
}

func resourceIKEPolicyDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    policyID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v3/%s/ike_policies/%s/", platform, policyID)
    _, err := config.doNormalRequest(platform, resourcePath, "DELETE", nil)

    if err != nil {
        return fmt.Errorf("Unable to delete IKE policy %s: on %s %v", policyID, platform, err)
    }

    stateConf := &resource.StateChangeConf{
        Pending:    []string{"DELETING"},
        Target:     []string{"DELETED"},
        Refresh:    IKEStateRefreshForDeletedFunc(config, platform, resourcePath),
        Timeout:    d.Timeout(schema.TimeoutDelete),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf( 
            "Error waiting for twcc_ike_policy %s to become DELETED: %v", policyID, err)
    }

    d.SetId("")

    return nil
}
