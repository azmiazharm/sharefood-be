// Package storage
package storage

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Compile time check to verify implements the Storage interface.
var _ Storage = (*awss3)(nil)

// awss3 implements the storage interface to provides the ability
// put get delete files to AWS S3.
type awss3 struct {
	svc *s3.S3
}

// NewAwsS3 creates a new instance AWS S3 Service,
func NewAwsS3(sess *session.Session) Storage {
	return &awss3{
		svc: s3.New(sess),
	}
}

// Put a new S3 object or overwrites an existing one.
func (s *awss3) Put(ctx context.Context, bucket, key string, contents []byte, cacheAble bool, contentType string) error {
	cacheControl := "public, max-age=86400"
	if !cacheAble {
		cacheControl = "no-cache, max-age=0"
	}

	putInput := s3.PutObjectInput{
		Bucket:       aws.String(bucket),
		Key:          aws.String(key),
		CacheControl: aws.String(cacheControl),
		Body:         bytes.NewReader(contents),
	}
	if contentType != "" {
		putInput.ContentType = aws.String(contentType)
	}

	if _, err := s.svc.PutObjectWithContext(ctx, &putInput); err != nil {
		return fmt.Errorf("storage create object: %w", err)
	}

	return nil
}

// Delete deletes a S3 object, returns nil if the object was successfully
// deleted, or of the object doesn't exist.
func (s *awss3) Delete(ctx context.Context, bucket, key string) error {
	_, err := s.svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("storage delete object: %w", err)
	}
	return nil
}

// Get returns the contents for the given object. If the object does not
// exist will returns error not found
func (s *awss3) Get(ctx context.Context, bucket, key string) ([]byte, error) {
	o, err := s.svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		if aErr, ok := err.(awserr.Error); ok &&
			aErr.Code() == s3.ErrCodeNoSuchBucket || aErr.Code() == s3.ErrCodeNoSuchKey {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	defer o.Body.Close()

	b, err := ioutil.ReadAll(o.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}

	return b, nil
}
