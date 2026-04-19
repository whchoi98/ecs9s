package aws

import (
	"context"
	"fmt"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsClient struct {
	client *sm.Client
}

func NewSecretsClient(cfg awssdk.Config) *SecretsClient {
	return &SecretsClient{client: sm.NewFromConfig(cfg)}
}

type Secret struct {
	Name            string
	ARN             string
	Description     string
	LastRotated     time.Time
	LastChanged     time.Time
	RotationEnabled bool
}

func (c *SecretsClient) ListSecrets(ctx context.Context) ([]Secret, error) {
	out, err := c.client.ListSecrets(ctx, &sm.ListSecretsInput{
		MaxResults: awssdk.Int32(100),
	})
	if err != nil {
		return nil, fmt.Errorf("list secrets: %w", err)
	}

	var secrets []Secret
	for _, s := range out.SecretList {
		var lastRotated, lastChanged time.Time
		if s.LastRotatedDate != nil {
			lastRotated = *s.LastRotatedDate
		}
		if s.LastChangedDate != nil {
			lastChanged = *s.LastChangedDate
		}

		rotEnabled := false
		if s.RotationEnabled != nil {
			rotEnabled = *s.RotationEnabled
		}

		secrets = append(secrets, Secret{
			Name:            awssdk.ToString(s.Name),
			ARN:             awssdk.ToString(s.ARN),
			Description:     awssdk.ToString(s.Description),
			LastRotated:     lastRotated,
			LastChanged:     lastChanged,
			RotationEnabled: rotEnabled,
		})
	}
	return secrets, nil
}
