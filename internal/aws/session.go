package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

type Session struct {
	Config  aws.Config
	Profile string
	Region  string
}

func NewSession(profile, region string) (*Session, error) {
	opts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(region),
	}
	if profile != "" && profile != "default" {
		opts = append(opts, awsconfig.WithSharedConfigProfile(profile))
	}

	cfg, err := awsconfig.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	return &Session{
		Config:  cfg,
		Profile: profile,
		Region:  region,
	}, nil
}

func (s *Session) SwitchProfile(profile string) error {
	ns, err := NewSession(profile, s.Region)
	if err != nil {
		return err
	}
	*s = *ns
	return nil
}

func (s *Session) SwitchRegion(region string) error {
	ns, err := NewSession(s.Profile, region)
	if err != nil {
		return err
	}
	*s = *ns
	return nil
}
