package aws

import (
	"context"
	"fmt"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SSMClient struct {
	client *ssm.Client
}

func NewSSMClient(cfg awssdk.Config) *SSMClient {
	return &SSMClient{client: ssm.NewFromConfig(cfg)}
}

type Parameter struct {
	Name         string
	Type         string
	Value        string // Always masked for SecureString; plaintext only for String/StringList
	Version      int64
	LastModified time.Time
}

const secureStringMask = "****"

func (c *SSMClient) ListParameters(ctx context.Context, prefix string) ([]Parameter, error) {
	// Never decrypt — SecureString values come back as encrypted blobs
	input := &ssm.GetParametersByPathInput{
		Path:           &prefix,
		Recursive:      awssdk.Bool(true),
		WithDecryption: awssdk.Bool(false),
		MaxResults:     awssdk.Int32(50),
	}

	out, err := c.client.GetParametersByPath(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("get parameters by path: %w", err)
	}

	var params []Parameter
	for _, p := range out.Parameters {
		value := awssdk.ToString(p.Value)
		// Mask SecureString values — they're returned as encrypted blobs when
		// WithDecryption is false, but we mask them explicitly for safety.
		if string(p.Type) == "SecureString" {
			value = secureStringMask
		}
		params = append(params, Parameter{
			Name:         awssdk.ToString(p.Name),
			Type:         string(p.Type),
			Value:        value,
			Version:      p.Version,
			LastModified: awssdk.ToTime(p.LastModifiedDate),
		})
	}
	return params, nil
}

func (c *SSMClient) GetParameter(ctx context.Context, name string) (*Parameter, error) {
	// Never decrypt SecureString values — read without decryption
	out, err := c.client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: awssdk.Bool(false),
	})
	if err != nil {
		return nil, fmt.Errorf("get parameter: %w", err)
	}

	value := awssdk.ToString(out.Parameter.Value)
	if string(out.Parameter.Type) == "SecureString" {
		value = secureStringMask
	}

	return &Parameter{
		Name:         awssdk.ToString(out.Parameter.Name),
		Type:         string(out.Parameter.Type),
		Value:        value,
		Version:      out.Parameter.Version,
		LastModified: awssdk.ToTime(out.Parameter.LastModifiedDate),
	}, nil
}
