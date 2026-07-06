package database

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetS3Client(logger *log.Logger) *s3.Client {
	ctx := context.TODO()
	cnf, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		logger.Fatalf("Unable to Load SDK Configs %s", err)
	}

	return s3.NewFromConfig(cnf)
}
