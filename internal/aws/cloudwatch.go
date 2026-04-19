package aws

import (
	"context"
	"fmt"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

type CloudWatchClient struct {
	logs    *cloudwatchlogs.Client
	metrics *cloudwatch.Client
}

func NewCloudWatchClient(cfg awssdk.Config) *CloudWatchClient {
	return &CloudWatchClient{
		logs:    cloudwatchlogs.NewFromConfig(cfg),
		metrics: cloudwatch.NewFromConfig(cfg),
	}
}

type LogEvent struct {
	Timestamp time.Time
	Message   string
}

func (c *CloudWatchClient) GetLogEvents(ctx context.Context, logGroup, logStream string, startTime time.Time, limit int32) ([]LogEvent, error) {
	input := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  &logGroup,
		LogStreamName: &logStream,
		StartFromHead: awssdk.Bool(false),
		Limit:         &limit,
	}
	if !startTime.IsZero() {
		ms := startTime.UnixMilli()
		input.StartTime = &ms
	}

	out, err := c.logs.GetLogEvents(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("get log events: %w", err)
	}

	var events []LogEvent
	for _, e := range out.Events {
		ts := time.UnixMilli(awssdk.ToInt64(e.Timestamp))
		events = append(events, LogEvent{
			Timestamp: ts,
			Message:   awssdk.ToString(e.Message),
		})
	}
	return events, nil
}

func (c *CloudWatchClient) ListLogGroups(ctx context.Context, prefix string) ([]string, error) {
	input := &cloudwatchlogs.DescribeLogGroupsInput{}
	if prefix != "" {
		input.LogGroupNamePrefix = &prefix
	}

	out, err := c.logs.DescribeLogGroups(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("describe log groups: %w", err)
	}

	var groups []string
	for _, g := range out.LogGroups {
		groups = append(groups, awssdk.ToString(g.LogGroupName))
	}
	return groups, nil
}

func (c *CloudWatchClient) ListLogStreams(ctx context.Context, logGroup string) ([]string, error) {
	out, err := c.logs.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: &logGroup,
		OrderBy:      "LastEventTime",
		Descending:   awssdk.Bool(true),
		Limit:        awssdk.Int32(20),
	})
	if err != nil {
		return nil, fmt.Errorf("describe log streams: %w", err)
	}

	var streams []string
	for _, s := range out.LogStreams {
		streams = append(streams, awssdk.ToString(s.LogStreamName))
	}
	return streams, nil
}

type MetricDatapoint struct {
	Timestamp time.Time
	Value     float64
}

func (c *CloudWatchClient) GetECSMetrics(ctx context.Context, clusterName, serviceName, metricName string, duration time.Duration) ([]MetricDatapoint, error) {
	end := time.Now()
	start := end.Add(-duration)
	period := int32(300) // 5 minutes

	out, err := c.metrics.GetMetricStatistics(ctx, &cloudwatch.GetMetricStatisticsInput{
		Namespace:  awssdk.String("AWS/ECS"),
		MetricName: &metricName,
		Dimensions: []cwtypes.Dimension{
			{Name: awssdk.String("ClusterName"), Value: &clusterName},
			{Name: awssdk.String("ServiceName"), Value: &serviceName},
		},
		StartTime:  &start,
		EndTime:    &end,
		Period:     &period,
		Statistics: []cwtypes.Statistic{cwtypes.StatisticAverage},
	})
	if err != nil {
		return nil, fmt.Errorf("get metrics: %w", err)
	}

	var points []MetricDatapoint
	for _, dp := range out.Datapoints {
		points = append(points, MetricDatapoint{
			Timestamp: awssdk.ToTime(dp.Timestamp),
			Value:     awssdk.ToFloat64(dp.Average),
		})
	}
	return points, nil
}

// --- Alarms ---

type Alarm struct {
	Name       string
	State      string
	MetricName string
	Namespace  string
	Threshold  float64
	Comparison string
	UpdatedAt  time.Time
}

func (c *CloudWatchClient) ListAlarms(ctx context.Context) ([]Alarm, error) {
	out, err := c.metrics.DescribeAlarms(ctx, &cloudwatch.DescribeAlarmsInput{
		MaxRecords: awssdk.Int32(100),
	})
	if err != nil {
		return nil, fmt.Errorf("describe alarms: %w", err)
	}

	var alarms []Alarm
	for _, a := range out.MetricAlarms {
		alarms = append(alarms, Alarm{
			Name:       awssdk.ToString(a.AlarmName),
			State:      string(a.StateValue),
			MetricName: awssdk.ToString(a.MetricName),
			Namespace:  awssdk.ToString(a.Namespace),
			Threshold:  awssdk.ToFloat64(a.Threshold),
			Comparison: string(a.ComparisonOperator),
			UpdatedAt:  awssdk.ToTime(a.StateUpdatedTimestamp),
		})
	}
	return alarms, nil
}
