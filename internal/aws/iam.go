package aws

import (
	"context"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

type IAMClient struct {
	client *iam.Client
}

func NewIAMClient(cfg awssdk.Config) *IAMClient {
	return &IAMClient{client: iam.NewFromConfig(cfg)}
}

type Role struct {
	Name       string
	ARN        string
	Path       string
	CreateDate string
	Policies   []string
}

func (c *IAMClient) ListRoles(ctx context.Context, pathPrefix string) ([]Role, error) {
	input := &iam.ListRolesInput{}
	if pathPrefix != "" {
		input.PathPrefix = &pathPrefix
	}

	out, err := c.client.ListRoles(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("list roles: %w", err)
	}

	var roles []Role
	for _, r := range out.Roles {
		roles = append(roles, Role{
			Name:       awssdk.ToString(r.RoleName),
			ARN:        awssdk.ToString(r.Arn),
			Path:       awssdk.ToString(r.Path),
			CreateDate: r.CreateDate.Format("2006-01-02"),
		})
	}
	return roles, nil
}

func (c *IAMClient) GetRolePolicies(ctx context.Context, roleName string) ([]string, error) {
	// Attached managed policies
	attachedOut, err := c.client.ListAttachedRolePolicies(ctx, &iam.ListAttachedRolePoliciesInput{
		RoleName: &roleName,
	})
	if err != nil {
		return nil, fmt.Errorf("list attached policies: %w", err)
	}

	var policies []string
	for _, p := range attachedOut.AttachedPolicies {
		policies = append(policies, awssdk.ToString(p.PolicyName))
	}

	// Inline policies
	inlineOut, err := c.client.ListRolePolicies(ctx, &iam.ListRolePoliciesInput{
		RoleName: &roleName,
	})
	if err == nil {
		for _, name := range inlineOut.PolicyNames {
			policies = append(policies, "(inline) "+name)
		}
	}

	return policies, nil
}
