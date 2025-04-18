package aws

import (
	"context"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

var (
	once sync.Once
	s    aws.Config
)

func NewAWSConfig() aws.Config {
	once.Do(func() {
		region := os.Getenv("AWS_REGION")
		if region == "" {
			region = "ap-northeast-1" // Default region
		}
		c, err := awsConfig.LoadDefaultConfig(context.TODO(),
			awsConfig.WithRegion(region),
			awsConfig.WithSharedConfigProfile("admin"),
		)
		if err != nil {
			return
		}

		s = c
	})

	return s
}

func NewSESClient() *sesv2.Client {
	sesClient := sesv2.NewFromConfig(NewAWSConfig())
	return sesClient
}
