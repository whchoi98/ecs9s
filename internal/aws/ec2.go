package aws

import (
	"context"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Client struct {
	client *ec2.Client
}

func NewEC2Client(cfg awssdk.Config) *EC2Client {
	return &EC2Client{client: ec2.NewFromConfig(cfg)}
}

type VPC struct {
	ID        string
	CIDR      string
	State     string
	IsDefault bool
	Name      string
}

type Subnet struct {
	ID        string
	VPCID     string
	CIDR      string
	AZ        string
	Available int32
	Name      string
}

type SecurityGroup struct {
	ID          string
	Name        string
	Description string
	VPCID       string
	IngressCount int
	EgressCount  int
}

type Instance struct {
	ID           string
	Type         string
	State        string
	PrivateIP    string
	PublicIP     string
	SubnetID     string
	AMI          string
	Name         string
}

func (c *EC2Client) ListVPCs(ctx context.Context) ([]VPC, error) {
	out, err := c.client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{})
	if err != nil {
		return nil, fmt.Errorf("describe vpcs: %w", err)
	}

	var vpcs []VPC
	for _, v := range out.Vpcs {
		name := tagName(v.Tags)
		vpcs = append(vpcs, VPC{
			ID:        awssdk.ToString(v.VpcId),
			CIDR:      awssdk.ToString(v.CidrBlock),
			State:     string(v.State),
			IsDefault: awssdk.ToBool(v.IsDefault),
			Name:      name,
		})
	}
	return vpcs, nil
}

func (c *EC2Client) ListSubnets(ctx context.Context, vpcID string) ([]Subnet, error) {
	input := &ec2.DescribeSubnetsInput{}
	if vpcID != "" {
		input.Filters = []ec2types.Filter{
			{Name: awssdk.String("vpc-id"), Values: []string{vpcID}},
		}
	}

	out, err := c.client.DescribeSubnets(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("describe subnets: %w", err)
	}

	var subnets []Subnet
	for _, s := range out.Subnets {
		subnets = append(subnets, Subnet{
			ID:        awssdk.ToString(s.SubnetId),
			VPCID:     awssdk.ToString(s.VpcId),
			CIDR:      awssdk.ToString(s.CidrBlock),
			AZ:        awssdk.ToString(s.AvailabilityZone),
			Available: awssdk.ToInt32(s.AvailableIpAddressCount),
			Name:      tagName(s.Tags),
		})
	}
	return subnets, nil
}

func (c *EC2Client) ListSecurityGroups(ctx context.Context, vpcID string) ([]SecurityGroup, error) {
	input := &ec2.DescribeSecurityGroupsInput{}
	if vpcID != "" {
		input.Filters = []ec2types.Filter{
			{Name: awssdk.String("vpc-id"), Values: []string{vpcID}},
		}
	}

	out, err := c.client.DescribeSecurityGroups(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("describe security groups: %w", err)
	}

	var sgs []SecurityGroup
	for _, sg := range out.SecurityGroups {
		sgs = append(sgs, SecurityGroup{
			ID:           awssdk.ToString(sg.GroupId),
			Name:         awssdk.ToString(sg.GroupName),
			Description:  awssdk.ToString(sg.Description),
			VPCID:        awssdk.ToString(sg.VpcId),
			IngressCount: len(sg.IpPermissions),
			EgressCount:  len(sg.IpPermissionsEgress),
		})
	}
	return sgs, nil
}

func (c *EC2Client) ListInstances(ctx context.Context) ([]Instance, error) {
	out, err := c.client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("describe instances: %w", err)
	}

	var instances []Instance
	for _, r := range out.Reservations {
		for _, i := range r.Instances {
			instances = append(instances, Instance{
				ID:        awssdk.ToString(i.InstanceId),
				Type:      string(i.InstanceType),
				State:     string(i.State.Name),
				PrivateIP: awssdk.ToString(i.PrivateIpAddress),
				PublicIP:  awssdk.ToString(i.PublicIpAddress),
				SubnetID:  awssdk.ToString(i.SubnetId),
				AMI:       awssdk.ToString(i.ImageId),
				Name:      tagName(i.Tags),
			})
		}
	}
	return instances, nil
}

func tagName(tags []ec2types.Tag) string {
	for _, t := range tags {
		if awssdk.ToString(t.Key) == "Name" {
			return awssdk.ToString(t.Value)
		}
	}
	return ""
}
