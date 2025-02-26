package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/talvor/asyncapi/config"
)

func main() {
	ctx := context.Background()
	conf := config.GetConfig()

	sdkConfig, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("Couldn't load default config.  Have you set up your AWS credentials?")
		fmt.Println(err)
		return
	}

	s3Client := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(conf.S3Endpoint)
		o.UsePathStyle = true
	})
	out, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("failed to list buckets, %v", err)
	}

	for _, bucket := range out.Buckets {
		fmt.Println(*bucket.Name)
	}

	sqsClient := sqs.NewFromConfig(sdkConfig, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String(conf.SQSEndpoint)
	})

	out2, err := sqsClient.ListQueues(ctx, &sqs.ListQueuesInput{})
	if err != nil {
		log.Fatalf("failed to list queues, %v", err)
	}

	for _, queue := range out2.QueueUrls {
		fmt.Println(queue)
	}
}
