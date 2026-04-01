package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	Client                   *s3.Client
	AddDisableEndpointPrefix func(*s3.Options)
}

func (s *S3Service) PutObject(ctx context.Context, input *s3.PutObjectInput, opt ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	input.Key = TrimKey(input.Key)
	return s.Client.PutObject(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) ListParts(ctx context.Context, input *s3.ListPartsInput, opt ...func(*s3.Options)) (*s3.ListPartsOutput, error) {
	input.Key = TrimKey(input.Key)
	return s.Client.ListParts(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) UploadPart(ctx context.Context, input *s3.UploadPartInput, opt ...func(*s3.Options)) (*s3.UploadPartOutput, error) {
	input.Key = TrimKey(input.Key)
	return s.Client.UploadPart(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) GetObject(ctx context.Context, input *s3.GetObjectInput, opt ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	input.Key = TrimKey(input.Key)
	return s.Client.GetObject(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) HeadObject(ctx context.Context, input *s3.HeadObjectInput, opt ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	input.Key = TrimKey(input.Key)
	return s.Client.HeadObject(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) CreateMultipartUpload(ctx context.Context, input *s3.CreateMultipartUploadInput, opt ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error) {
	input.Key = TrimKey(input.Key)
	res, err := s.Client.CreateMultipartUpload(ctx, input, s.AddDisableEndpointPrefix)
	if err != nil {
		return nil, fmt.Errorf("couuld not create multipart upload: %w", err)
	}
	return res, nil
}
func (s *S3Service) AbortMultipartUpload(ctx context.Context, input *s3.AbortMultipartUploadInput, opt ...func(*s3.Options)) (*s3.AbortMultipartUploadOutput, error) {
	input.Key = TrimKey(input.Key)
	return s.Client.AbortMultipartUpload(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, opt ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	input.Key = TrimKey(input.Key)
	return s.Client.DeleteObject(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) DeleteObjects(ctx context.Context, input *s3.DeleteObjectsInput, opt ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error) {
	return s.Client.DeleteObjects(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) CompleteMultipartUpload(ctx context.Context, input *s3.CompleteMultipartUploadInput, opt ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error) {
	input.Key = TrimKey(input.Key)
	return s.Client.CompleteMultipartUpload(ctx, input, s.AddDisableEndpointPrefix)
}
func (s *S3Service) UploadPartCopy(ctx context.Context, input *s3.UploadPartCopyInput, opt ...func(*s3.Options)) (*s3.UploadPartCopyOutput, error) {
	input.Key = TrimKey(input.Key)
	return s.Client.UploadPartCopy(ctx, input, s.AddDisableEndpointPrefix)
}

var trimKeyTenantPrefixRe = regexp.MustCompile(
	`^tenantname-([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})-(.+)$`,
)

func TrimKey(key *string) *string {
	if key == nil || *key == "" {
		return key
	}

	parts := strings.Split(*key, "/")
	if len(parts) < 2 {
		return key
	}

	prefix := parts[0]
	last := parts[len(parts)-1]

	m := trimKeyTenantPrefixRe.FindStringSubmatch(last)
	if len(m) != 3 {
		return key
	}

	trimmedFile := m[2]
	res := prefix + "/" + trimmedFile
	return &res

}
