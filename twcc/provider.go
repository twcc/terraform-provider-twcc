package twcc
  
import (
    "github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var osMutexKV = mutexkv.NewMutexKV()

type PConfig struct {
    Config
}

// Provider returns a schema.Provider for TWCC.
func Provider() terraform.ResourceProvider {
    provider := &schema.Provider{
        Schema: map[string]*schema.Schema{
            "apikey": {
                Type:		schema.TypeString,
                Required:	true,
                DefaultFunc:	schema.EnvDefaultFunc("TWCC_APIKEY", ""),
                Description:	descriptions["apikey"],
            },
            "apigw_url": {       
                Type:		schema.TypeString,
                Required:	true,
                DefaultFunc:	schema.EnvDefaultFunc("APIGW_URL", ""),
                Description:	descriptions["apigw_url"],
            },
        },

        DataSourcesMap: map[string]*schema.Resource{
            "twcc_network":			dataSourceNetwork(),
            "twcc_project":			dataSourceProject(),
            "twcc_solution":			dataSourceSolution(),
            "twcc_firewall":			dataSourceFirewall(),
            "twcc_firewall_rule":		dataSourceFirewallRule(),
            "twcc_vcs":				dataSourceVCS(),
            "twcc_volume":			dataSourceVolume(),
            "twcc_volume_snapshot":		dataSourceVolumeSnapshot(),
            "twcc_ike_policy":			dataSourceIKEPolicy(),
            "twcc_ipsec_policy":		dataSourceIPSecPolicy(),
            "twcc_vpn":				dataSourceVPN(),
            "twcc_container":			dataSourceContainer(),
            "twcc_s3_key":			dataSourceS3Key(),
            "twcc_secret":			dataSourceSecret(),
            "twcc_security_group":		dataSourceSecurityGroup(),
            "twcc_loadbalancer":		dataSourceLoadBalancer(),
            "twcc_auto_scaling_policy":		dataSourceAutoScalingPolicy(),
            "twcc_extra_property":		dataSourceExtraProperty(),
            "twcc_waf":				dataSourceWAF(),
        },

        ResourcesMap: map[string]*schema.Resource{
            "twcc_auto_scaling_policy":		resourceAutoScalingPolicy(),
            "twcc_auto_scaling_relation":	resourceAutoScalingRelation(),
            "twcc_container":			resourceContainer(),
            "twcc_firewall":			resourceFirewall(),
            "twcc_firewall_rule":		resourceFirewallRule(),
            "twcc_ike_policy":			resourceIKEPolicy(),
            "twcc_ipsec_policy":		resourceIPSecPolicy(),
            "twcc_loadbalancer":		resourceLoadBalancer(),
            "twcc_network":			resourceNetwork(),
            "twcc_vcs":				resourceVCS(),
            "twcc_vcs_image":			resourceVCSImage(),
            "twcc_volume":			resourceVolume(),
            "twcc_volume_attachment":		resourceVolumeAttachment(),
            "twcc_volume_snapshot":		resourceVolumeSnapshot(),
            "twcc_vpn":				resourceVPN(),
            "twcc_vpn_connection":		resourceVPNConnection(),
            "twcc_s3_key":			resourceS3Key(),
            "twcc_secret":			resourceSecret(),
            "twcc_security_group_rule":		resourceSecurityGroupRule(),
            "twcc_waf":				resourceWAF(),
        },
    }

    provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
        terraformVersion := provider.TerraformVersion
        if terraformVersion == "" {
            // Terraform 0.12 introduced this field to the protocol
            // We can therefore assume that if it's missing it's 0.10 or 0.11
            terraformVersion = "0.11+compatible"
        }
        return configureProvider(d, terraformVersion)
    }
    return provider
}

var descriptions map[string]string

func init() {
    descriptions = map[string]string{
        "apikey": "APIKey to login with.",
        "apigw_url": "APIGW endpoint to request to.",
    }
}

func configureProvider(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
    config := PConfig{
        Config{
            TWCC_APIKEY:	d.Get("apikey").(string),
            APIGW_URL:		d.Get("apigw_url").(string),
        },
    }

    if err := config.LoadAndValidate(); err != nil {
        return nil, err
    }

    return &config, nil
}
