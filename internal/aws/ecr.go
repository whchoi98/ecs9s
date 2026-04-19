package aws

import (
	"context"
	"fmt"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

type ECRClient struct {
	client *ecr.Client
}

func NewECRClient(cfg awssdk.Config) *ECRClient {
	return &ECRClient{client: ecr.NewFromConfig(cfg)}
}

type Repository struct {
	Name      string
	URI       string
	CreatedAt time.Time
	ImageCount int
}

type Image struct {
	Tags     string
	Digest   string
	SizeMB   float64
	PushedAt time.Time
}

func (c *ECRClient) ListRepositories(ctx context.Context) ([]Repository, error) {
	out, err := c.client.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{})
	if err != nil {
		return nil, fmt.Errorf("describe repositories: %w", err)
	}

	var repos []Repository
	for _, r := range out.Repositories {
		imgOut, _ := c.client.DescribeImages(ctx, &ecr.DescribeImagesInput{
			RepositoryName: r.RepositoryName,
		})
		imgCount := 0
		if imgOut != nil {
			imgCount = len(imgOut.ImageDetails)
		}

		repos = append(repos, Repository{
			Name:       awssdk.ToString(r.RepositoryName),
			URI:        awssdk.ToString(r.RepositoryUri),
			CreatedAt:  awssdk.ToTime(r.CreatedAt),
			ImageCount: imgCount,
		})
	}
	return repos, nil
}

func (c *ECRClient) ListImages(ctx context.Context, repoName string) ([]Image, error) {
	out, err := c.client.DescribeImages(ctx, &ecr.DescribeImagesInput{
		RepositoryName: &repoName,
	})
	if err != nil {
		return nil, fmt.Errorf("describe images: %w", err)
	}

	var images []Image
	for _, img := range out.ImageDetails {
		var tags string
		for i, t := range img.ImageTags {
			if i > 0 {
				tags += ", "
			}
			tags += t
		}
		if tags == "" {
			tags = "<untagged>"
		}

		sizeMB := float64(awssdk.ToInt64(img.ImageSizeInBytes)) / 1024 / 1024
		images = append(images, Image{
			Tags:     tags,
			Digest:   awssdk.ToString(img.ImageDigest),
			SizeMB:   sizeMB,
			PushedAt: awssdk.ToTime(img.ImagePushedAt),
		})
	}
	return images, nil
}
