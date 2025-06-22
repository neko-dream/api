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
		region := "ap-northeast-1"
		var funs []func(*awsConfig.LoadOptions) error
		// localならWithSharedConfigProfileを使う
		if os.Getenv("ENV") == "local" {
			funs = append(funs, awsConfig.WithSharedConfigProfile("local"))
		}
		funs = append(funs, awsConfig.WithRegion(region))

		c, err := awsConfig.LoadDefaultConfig(context.TODO(), funs...)
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
