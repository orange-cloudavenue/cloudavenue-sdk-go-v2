/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package s3

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Client provides CloudAvenue S3 operations.
type Client struct {
	client *s3.Client
}

// NewClient creates a new S3 client from an aws.Config.
func NewClient(cfg aws.Config) *Client {
	return &Client{
		client: s3.NewFromConfig(cfg),
	}
}

// ListBuckets lists all buckets in the S3 account.
func (c *Client) ListBuckets(ctx context.Context) ([]types.Bucket, error) {
	resp, err := c.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	return resp.Buckets, nil
}

// ListObjects lists objects in a bucket with optional prefix.
func (c *Client) ListObjects(ctx context.Context, bucket, prefix string) ([]types.Object, error) {
	var objects []types.Object
	paginator := s3.NewListObjectsV2Paginator(c.client, &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		objects = append(objects, page.Contents...)
	}

	return objects, nil
}

// PutObject uploads an object to a bucket.
func (c *Client) PutObject(ctx context.Context, bucket, key string, body io.Reader, contentType string) (*s3.PutObjectOutput, error) {
	return c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		Body:        body,
		ContentType: &contentType,
	})
}

// GetObject downloads an object from a bucket.
func (c *Client) GetObject(ctx context.Context, bucket, key string) (*s3.GetObjectOutput, error) {
	return c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
}

// DeleteObject deletes an object from a bucket.
func (c *Client) DeleteObject(ctx context.Context, bucket, key string) (*s3.DeleteObjectOutput, error) {
	return c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
}

// CreateBucket creates a new bucket.
func (c *Client) CreateBucket(ctx context.Context, bucket string) (*s3.CreateBucketOutput, error) {
	return c.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucket,
	})
}

// DeleteBucket deletes a bucket.
func (c *Client) DeleteBucket(ctx context.Context, bucket string) (*s3.DeleteBucketOutput, error) {
	return c.client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: &bucket,
	})
}

// HeadBucket checks if a bucket exists and is accessible.
func (c *Client) HeadBucket(ctx context.Context, bucket string) (*s3.HeadBucketOutput, error) {
	return c.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &bucket,
	})
}

// GetObjectPresignedURL generates a presigned URL for getting an object.
func (c *Client) GetObjectPresignedURL(ctx context.Context, bucket, key string, expires time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(c.client)
	resp, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}, s3.WithPresignExpires(expires))
	if err != nil {
		return "", err
	}
	return resp.URL, nil
}
