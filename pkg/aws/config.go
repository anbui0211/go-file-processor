package paws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

const (
	awsEndpoint = "http://localhost:4566"
	awsRegion   = "us-east-1"
)

func LoadAWSConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		log.Fatalf("Cannot load AWS configs: %s", err)
	}
	return cfg
}
