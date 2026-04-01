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

var trimKeyUUIDRe = regexp.MustCompile(
	`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`,
)

func TrimKey(key *string) *string {
	if key == nil || *key == "" {
		return key
	}

	parts := strings.Split(*key, "/")
	if len(parts) < 2 {
		return key
	}

	// Keep only the first path segment, regardless of how deep the input key is.
	prefix := parts[0]

	// Extract filename from the last path segment: everything after "<uuid>-"
	last := parts[len(parts)-1]
	loc := trimKeyUUIDRe.FindStringIndex(last)
	if loc == nil {
		return key
	}

	uuidEnd := loc[1]
	if uuidEnd >= len(last) || last[uuidEnd] != '-' || uuidEnd+1 >= len(last) {
		return key
	}

	filename := last[uuidEnd+1:]
	res := prefix + "/" + filename
	return &res

}
