package depoplambda

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"aws_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "us-east-1",
				Description: "AWS Region",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"depoplambda_lambdas": dataSourceLambdas(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		return providerConfigure(d, provider.TerraformVersion)
	}

	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		AWSRegion:        d.Get("aws_region").(string),
		terraformVersion: terraformVersion,
	}
	return config.Client()
}
