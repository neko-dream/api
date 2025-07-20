package repository

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/neko-dream/server/internal/infrastructure/config"
)

var (
	once sync.Once
	conf aws.Config
)

// initConfig awsConfigを作成。 otelawsによる計装も設定
func InitConfig(appConf *config.Config) aws.Config {
	return InitConfigWithContext(context.Background(), appConf)
}

func InitConfigWithContext(ctx context.Context, appConf *config.Config) aws.Config {
	once.Do(func() {
		c, err := awsConfig.LoadDefaultConfig(
			ctx,
			awsConfig.WithRegion(appConf.R2_REGION),
			awsConfig.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(
					appConf.R2_ACCESS_KEY_ID,
					appConf.R2_SECRET_ACCESS_KEY,
					"",
				),
			),
		)
		if err != nil {
			return
		}

		conf = c
	})

	return conf
}

func InitS3Client(appConf *config.Config) *s3.Client {
	client := s3.NewFromConfig(InitConfig(appConf), func(o *s3.Options) {
		o.BaseEndpoint = aws.String(appConf.AWS_S3_ENDPOINT)
	})
	return client
}
