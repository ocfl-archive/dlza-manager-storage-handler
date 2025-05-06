package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"strings"
)

type S3Service struct {
	Client                   *s3.Client
	AddDisableEndpointPrefix func(*s3.Options)
}

func (s *S3Service) PutObject(ctx context.Context, input *s3.PutObjectInput, opt ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	input.Key = trimKey(input.Key)
	return s.Client.PutObject(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) ListParts(ctx context.Context, input *s3.ListPartsInput, opt ...func(*s3.Options)) (*s3.ListPartsOutput, error) {
	input.Key = trimKey(input.Key)
	return s.Client.ListParts(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) UploadPart(ctx context.Context, input *s3.UploadPartInput, opt ...func(*s3.Options)) (*s3.UploadPartOutput, error) {
	input.Key = trimKey(input.Key)
	return s.Client.UploadPart(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) GetObject(ctx context.Context, input *s3.GetObjectInput, opt ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	input.Key = trimKey(input.Key)
	return s.Client.GetObject(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) HeadObject(ctx context.Context, input *s3.HeadObjectInput, opt ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	input.Key = trimKey(input.Key)
	return s.Client.HeadObject(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) CreateMultipartUpload(ctx context.Context, input *s3.CreateMultipartUploadInput, opt ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error) {
	input.Key = trimKey(input.Key)
	res, err := s.Client.CreateMultipartUpload(ctx, input, s.AddDisableEndpointPrefix)
	if err != nil {
		return nil, fmt.Errorf("couuld not create multipart upload: %w", err)
	}
	return res, nil
}
func (s *S3Service) AbortMultipartUpload(ctx context.Context, input *s3.AbortMultipartUploadInput, opt ...func(*s3.Options)) (*s3.AbortMultipartUploadOutput, error) {
	input.Key = trimKey(input.Key)
	return s.Client.AbortMultipartUpload(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, opt ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	input.Key = trimKey(input.Key)
	return s.Client.DeleteObject(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) DeleteObjects(ctx context.Context, input *s3.DeleteObjectsInput, opt ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error) {
	return s.Client.DeleteObjects(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) CompleteMultipartUpload(ctx context.Context, input *s3.CompleteMultipartUploadInput, opt ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error) {
	input.Key = trimKey(input.Key)
	return s.Client.CompleteMultipartUpload(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) UploadPartCopy(ctx context.Context, input *s3.UploadPartCopyInput, opt ...func(*s3.Options)) (*s3.UploadPartCopyOutput, error) {
	input.Key = trimKey(input.Key)
	return s.Client.UploadPartCopy(ctx, input, s.AddDisableEndpointPrefix)
}

func trimKey(key *string) *string {
	formattedKeyArray := strings.Split(*key, "/")
	if len(formattedKeyArray) == 1 {
		return key
	}
	return &formattedKeyArray[1]
}
