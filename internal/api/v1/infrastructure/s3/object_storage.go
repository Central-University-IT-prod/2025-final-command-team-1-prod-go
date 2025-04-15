package s3_storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"time"

	config "example.com/m/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

var (
	S3Client ClientS3 = ClientS3{}
)

type ClientS3 struct {
	S3Client *s3.Client
}

type resolverV2 struct {
	BaseEndpoint string
}

func (basics ClientS3) UploadFile(ctx context.Context, bucketName string, objectKey string, file multipart.File) error {
	_, err := basics.S3Client.PutObject(ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			ChecksumSHA256: aws.String("UNSIGNED-PAYLOAD"),
			Body:   file,
			ACL:    "public-read",
			
		})
	if err != nil {
		return err
	} else {
		err = s3.NewObjectExistsWaiter(basics.S3Client).Wait(
			ctx, &s3.HeadObjectInput{Bucket: aws.String(bucketName), Key: aws.String(objectKey)}, time.Minute)
		if err != nil {
			return err
		}
	}
	return err
}

func (r *resolverV2) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	u, err := url.Parse(r.BaseEndpoint)
	if err != nil {
		return smithyendpoints.Endpoint{}, err
	}
	return smithyendpoints.Endpoint{URI: *u}, nil
}

func InitS3() {
	cfg, err := s3config.LoadDefaultConfig(
		context.TODO(),
		s3config.WithRegion(config.Config.S3Region),
		s3config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.Config.S3AWSAccessKeyID,
				config.Config.S3AwsSecretAccessKey,
				"",
			)),
	)
	if err != nil {
		fmt.Println(fmt.Printf("unable to load SDK config, %v", err))
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.EndpointResolverV2 = &resolverV2{BaseEndpoint: config.Config.S3Endpoint}
		o.UsePathStyle = true
	})
	S3Client = ClientS3{S3Client: client}
	fmt.Println(S3Client)
}
