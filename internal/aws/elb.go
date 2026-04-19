package aws

import (
	"context"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

type ELBClient struct {
	client *elbv2.Client
}

func NewELBClient(cfg awssdk.Config) *ELBClient {
	return &ELBClient{client: elbv2.NewFromConfig(cfg)}
}

type LoadBalancer struct {
	Name     string
	ARN      string
	DNSName  string
	Type     string
	Scheme   string
	State    string
	VPCID    string
}

type TargetGroup struct {
	Name         string
	ARN          string
	Protocol     string
	Port         int32
	HealthyCount int
	UnhealthyCount int
	TargetType   string
}

func (c *ELBClient) ListLoadBalancers(ctx context.Context) ([]LoadBalancer, error) {
	out, err := c.client.DescribeLoadBalancers(ctx, &elbv2.DescribeLoadBalancersInput{})
	if err != nil {
		return nil, fmt.Errorf("describe load balancers: %w", err)
	}

	var lbs []LoadBalancer
	for _, lb := range out.LoadBalancers {
		state := "unknown"
		if lb.State != nil {
			state = string(lb.State.Code)
		}
		lbs = append(lbs, LoadBalancer{
			Name:    awssdk.ToString(lb.LoadBalancerName),
			ARN:     awssdk.ToString(lb.LoadBalancerArn),
			DNSName: awssdk.ToString(lb.DNSName),
			Type:    string(lb.Type),
			Scheme:  string(lb.Scheme),
			State:   state,
			VPCID:   awssdk.ToString(lb.VpcId),
		})
	}
	return lbs, nil
}

func (c *ELBClient) ListTargetGroups(ctx context.Context) ([]TargetGroup, error) {
	out, err := c.client.DescribeTargetGroups(ctx, &elbv2.DescribeTargetGroupsInput{})
	if err != nil {
		return nil, fmt.Errorf("describe target groups: %w", err)
	}

	var tgs []TargetGroup
	for _, tg := range out.TargetGroups {
		healthy, unhealthy := 0, 0
		thOut, err := c.client.DescribeTargetHealth(ctx, &elbv2.DescribeTargetHealthInput{
			TargetGroupArn: tg.TargetGroupArn,
		})
		if err == nil {
			for _, th := range thOut.TargetHealthDescriptions {
				if th.TargetHealth != nil && string(th.TargetHealth.State) == "healthy" {
					healthy++
				} else {
					unhealthy++
				}
			}
		}

		tgs = append(tgs, TargetGroup{
			Name:           awssdk.ToString(tg.TargetGroupName),
			ARN:            awssdk.ToString(tg.TargetGroupArn),
			Protocol:       string(tg.Protocol),
			Port:           awssdk.ToInt32(tg.Port),
			HealthyCount:   healthy,
			UnhealthyCount: unhealthy,
			TargetType:     string(tg.TargetType),
		})
	}
	return tgs, nil
}
