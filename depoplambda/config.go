package depoplambda

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Config struct {
	AWSRegion        string
	terraformVersion string
}

type ManifestClient struct {
	s3conn     *s3.S3
	bucketName string
	bucketKey  string
}

func (c *Config) Client() (interface{}, error) {
	var client ManifestClient

	awsConfig := &aws.Config{
		Region: aws.String(c.AWSRegion),
	}

	awsSession, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}

	client.s3conn = s3.New(awsSession)

	return &client, nil
}
