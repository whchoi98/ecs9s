package aws

import (
	"context"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	appas "github.com/aws/aws-sdk-go-v2/service/applicationautoscaling"
	astypes "github.com/aws/aws-sdk-go-v2/service/applicationautoscaling/types"
)

type AutoScalingClient struct {
	client *appas.Client
}

func NewAutoScalingClient(cfg awssdk.Config) *AutoScalingClient {
	return &AutoScalingClient{client: appas.NewFromConfig(cfg)}
}

type ScalableTarget struct {
	ResourceID    string
	ServiceNS     string
	ScalableDim   string
	MinCapacity   int32
	MaxCapacity   int32
}

type ScalingPolicy struct {
	PolicyName string
	PolicyType string
	ResourceID string
	MetricName string
	TargetVal  float64
}

func (c *AutoScalingClient) ListScalableTargets(ctx context.Context) ([]ScalableTarget, error) {
	out, err := c.client.DescribeScalableTargets(ctx, &appas.DescribeScalableTargetsInput{
		ServiceNamespace: astypes.ServiceNamespaceEcs,
	})
	if err != nil {
		return nil, fmt.Errorf("describe scalable targets: %w", err)
	}

	var targets []ScalableTarget
	for _, t := range out.ScalableTargets {
		targets = append(targets, ScalableTarget{
			ResourceID:  awssdk.ToString(t.ResourceId),
			ServiceNS:   string(t.ServiceNamespace),
			ScalableDim: string(t.ScalableDimension),
			MinCapacity: awssdk.ToInt32(t.MinCapacity),
			MaxCapacity: awssdk.ToInt32(t.MaxCapacity),
		})
	}
	return targets, nil
}

func (c *AutoScalingClient) ListScalingPolicies(ctx context.Context) ([]ScalingPolicy, error) {
	out, err := c.client.DescribeScalingPolicies(ctx, &appas.DescribeScalingPoliciesInput{
		ServiceNamespace: astypes.ServiceNamespaceEcs,
	})
	if err != nil {
		return nil, fmt.Errorf("describe scaling policies: %w", err)
	}

	var policies []ScalingPolicy
	for _, p := range out.ScalingPolicies {
		var metricName string
		var targetVal float64
		if p.TargetTrackingScalingPolicyConfiguration != nil {
			targetVal = awssdk.ToFloat64(p.TargetTrackingScalingPolicyConfiguration.TargetValue)
			if p.TargetTrackingScalingPolicyConfiguration.PredefinedMetricSpecification != nil {
				metricName = string(p.TargetTrackingScalingPolicyConfiguration.PredefinedMetricSpecification.PredefinedMetricType)
			}
		}

		policies = append(policies, ScalingPolicy{
			PolicyName: awssdk.ToString(p.PolicyName),
			PolicyType: string(p.PolicyType),
			ResourceID: awssdk.ToString(p.ResourceId),
			MetricName: metricName,
			TargetVal:  targetVal,
		})
	}
	return policies, nil
}

func (c *AutoScalingClient) UpdateScalableTarget(ctx context.Context, resourceID string, min, max int32) error {
	_, err := c.client.RegisterScalableTarget(ctx, &appas.RegisterScalableTargetInput{
		ServiceNamespace:  astypes.ServiceNamespaceEcs,
		ResourceId:        &resourceID,
		ScalableDimension: astypes.ScalableDimensionECSServiceDesiredCount,
		MinCapacity:       &min,
		MaxCapacity:       &max,
	})
	return err
}
