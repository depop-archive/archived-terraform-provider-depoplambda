package depoplambda

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceLambdas() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLambdasRead,

		Schema: map[string]*schema.Schema{
			"s3_bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"s3_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "lambdas_manifest.json",
			},
			"manifests": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceLambdasRead(d *schema.ResourceData, meta interface{}) error {
	s3Client := meta.(*ManifestClient).s3conn

	headObjInput := &s3.HeadObjectInput{
		Bucket: aws.String(d.Get("s3_bucket").(string)),
		Key:    aws.String(d.Get("s3_key").(string)),
	}

	headObjOutput, err := s3Client.HeadObject(headObjInput)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
				return err
			default:
				return err
			}
		} else {
			return err
		}
	}

	d.SetId(aws.StringValue(headObjOutput.ETag))

	downloader := s3manager.NewDownloaderWithClient(s3Client)

	var buf []byte
	result := aws.NewWriteAtBuffer(buf)

	downloadObjInput := &s3.GetObjectInput{
		Bucket: aws.String(d.Get("s3_bucket").(string)),
		Key:    aws.String(d.Get("s3_key").(string)),
	}

	bytes, err := downloader.Download(result, downloadObjInput)

	if err != nil {
		return err
	}

	if bytes == 0 {
		return fmt.Errorf("downloaded file contained no data")
	}

	manifests := make([]map[string]interface{}, 0)
	err = json.Unmarshal(result.Bytes(), &manifests)
	if err != nil {
		return fmt.Errorf("Unable to decode manfiest files: %s", err)
	}

	if err := d.Set("manifests", manifests); err != nil {
		return fmt.Errorf("Error setting lambda manifests: %s", err)
	}

	return nil
}
