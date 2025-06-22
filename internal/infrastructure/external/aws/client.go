package aws

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

var (
	s         aws.Config
	once      sync.Once
	configErr error
)

func NewAWSConfig() (aws.Config, error) {
	once.Do(func() {
		region := "ap-northeast-1"
		var funs []func(*awsConfig.LoadOptions) error

		// localならWithSharedConfigProfileを使う
		if os.Getenv("ENV") == "local" {
			funs = append(funs, awsConfig.WithSharedConfigProfile("admin"))
		}
		funs = append(funs, awsConfig.WithRegion(region))

		var err error
		s, err = awsConfig.LoadDefaultConfig(context.TODO(), funs...)
		if err != nil {
			configErr = fmt.Errorf("failed to load AWS config: %w", err)
			return
		}
	})

	if configErr != nil {
		return aws.Config{}, configErr
	}

	return s, nil
}

func NewSESClient() *sesv2.Client {
	cfg, err := NewAWSConfig()
	if err != nil {
		fmt.Printf("Error creating AWS config: %v\n", err)
		return nil
	}
	// SESv2クライアントを作成
	sesClient := sesv2.NewFromConfig(cfg)
	if sesClient == nil {
		fmt.Println("Error creating SES client")
		return nil
	}
	return sesClient
}
