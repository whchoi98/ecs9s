package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type ECSClient struct {
	client *ecs.Client
}

func NewECSClient(cfg aws.Config) *ECSClient {
	return &ECSClient{client: ecs.NewFromConfig(cfg)}
}

// --- Data types ---

type Cluster struct {
	Name               string
	ARN                string
	Status             string
	ActiveServicesCount int32
	RunningTasksCount  int32
	RegisteredInstances int32
}

type Service struct {
	Name           string
	ARN            string
	Status         string
	DesiredCount   int32
	RunningCount   int32
	TaskDefinition string
	LaunchType     string
	LoadBalancers  string
	CreatedAt      time.Time
}

type Task struct {
	TaskARN        string
	TaskID         string
	Status         string
	TaskDefinition string
	LaunchType     string
	StartedAt      time.Time
	ContainerCount int
	PrivateIP      string
	Group          string
}

type Container struct {
	Name       string
	Image      string
	Status     string
	Health     string
	Ports      string
	RuntimeID  string
	TaskARN    string
	ClusterARN string
}

type TaskDefinition struct {
	Family        string
	Revision      int32
	ARN           string
	Status        string
	CPU           string
	Memory        string
	Compatibility string
	RegisteredAt  time.Time
}

// --- API Calls ---

func (c *ECSClient) ListClusters(ctx context.Context) ([]Cluster, error) {
	listOut, err := c.client.ListClusters(ctx, &ecs.ListClustersInput{})
	if err != nil {
		return nil, fmt.Errorf("list clusters: %w", err)
	}
	if len(listOut.ClusterArns) == 0 {
		return nil, nil
	}

	descOut, err := c.client.DescribeClusters(ctx, &ecs.DescribeClustersInput{
		Clusters: listOut.ClusterArns,
		Include:  []ecstypes.ClusterField{ecstypes.ClusterFieldStatistics},
	})
	if err != nil {
		return nil, fmt.Errorf("describe clusters: %w", err)
	}

	var clusters []Cluster
	for _, cl := range descOut.Clusters {
		clusters = append(clusters, Cluster{
			Name:                aws.ToString(cl.ClusterName),
			ARN:                 aws.ToString(cl.ClusterArn),
			Status:              aws.ToString(cl.Status),
			ActiveServicesCount: cl.ActiveServicesCount,
			RunningTasksCount:   cl.RunningTasksCount,
			RegisteredInstances: cl.RegisteredContainerInstancesCount,
		})
	}
	return clusters, nil
}

func (c *ECSClient) ListServices(ctx context.Context, clusterARN string) ([]Service, error) {
	listOut, err := c.client.ListServices(ctx, &ecs.ListServicesInput{
		Cluster: &clusterARN,
	})
	if err != nil {
		return nil, fmt.Errorf("list services: %w", err)
	}
	if len(listOut.ServiceArns) == 0 {
		return nil, nil
	}

	descOut, err := c.client.DescribeServices(ctx, &ecs.DescribeServicesInput{
		Cluster:  &clusterARN,
		Services: listOut.ServiceArns,
	})
	if err != nil {
		return nil, fmt.Errorf("describe services: %w", err)
	}

	var services []Service
	for _, svc := range descOut.Services {
		var lbs string
		for _, lb := range svc.LoadBalancers {
			if lbs != "" {
				lbs += ", "
			}
			lbs += aws.ToString(lb.TargetGroupArn)
		}
		services = append(services, Service{
			Name:           aws.ToString(svc.ServiceName),
			ARN:            aws.ToString(svc.ServiceArn),
			Status:         aws.ToString(svc.Status),
			DesiredCount:   svc.DesiredCount,
			RunningCount:   svc.RunningCount,
			TaskDefinition: shortTaskDef(aws.ToString(svc.TaskDefinition)),
			LaunchType:     string(svc.LaunchType),
			LoadBalancers:  lbs,
			CreatedAt:      aws.ToTime(svc.CreatedAt),
		})
	}
	return services, nil
}

func (c *ECSClient) ListTasks(ctx context.Context, clusterARN, serviceName string) ([]Task, error) {
	input := &ecs.ListTasksInput{Cluster: &clusterARN}
	if serviceName != "" {
		input.ServiceName = &serviceName
	}

	listOut, err := c.client.ListTasks(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	if len(listOut.TaskArns) == 0 {
		return nil, nil
	}

	descOut, err := c.client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: &clusterARN,
		Tasks:   listOut.TaskArns,
	})
	if err != nil {
		return nil, fmt.Errorf("describe tasks: %w", err)
	}

	var tasks []Task
	for _, t := range descOut.Tasks {
		var ip string
		for _, att := range t.Attachments {
			for _, kv := range att.Details {
				if aws.ToString(kv.Name) == "privateIPv4Address" {
					ip = aws.ToString(kv.Value)
				}
			}
		}

		arn := aws.ToString(t.TaskArn)
		tasks = append(tasks, Task{
			TaskARN:        arn,
			TaskID:         shortTaskID(arn),
			Status:         aws.ToString(t.LastStatus),
			TaskDefinition: shortTaskDef(aws.ToString(t.TaskDefinitionArn)),
			LaunchType:     string(t.LaunchType),
			StartedAt:      aws.ToTime(t.StartedAt),
			ContainerCount: len(t.Containers),
			PrivateIP:      ip,
			Group:          aws.ToString(t.Group),
		})
	}
	return tasks, nil
}

func (c *ECSClient) DescribeContainers(ctx context.Context, clusterARN, taskARN string) ([]Container, error) {
	descOut, err := c.client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: &clusterARN,
		Tasks:   []string{taskARN},
	})
	if err != nil {
		return nil, fmt.Errorf("describe task containers: %w", err)
	}
	if len(descOut.Tasks) == 0 {
		return nil, nil
	}

	var containers []Container
	for _, ctr := range descOut.Tasks[0].Containers {
		var ports string
		for _, p := range ctr.NetworkBindings {
			if ports != "" {
				ports += ", "
			}
			ports += fmt.Sprintf("%d->%d", aws.ToInt32(p.HostPort), aws.ToInt32(p.ContainerPort))
		}
		if ports == "" {
			for _, ni := range ctr.NetworkInterfaces {
				_ = ni
			}
		}

		health := "UNKNOWN"
		if ctr.HealthStatus != "" {
			health = string(ctr.HealthStatus)
		}

		containers = append(containers, Container{
			Name:       aws.ToString(ctr.Name),
			Image:      aws.ToString(ctr.Image),
			Status:     aws.ToString(ctr.LastStatus),
			Health:     health,
			Ports:      ports,
			RuntimeID:  aws.ToString(ctr.RuntimeId),
			TaskARN:    taskARN,
			ClusterARN: clusterARN,
		})
	}
	return containers, nil
}

func (c *ECSClient) ListTaskDefinitions(ctx context.Context) ([]TaskDefinition, error) {
	listOut, err := c.client.ListTaskDefinitionFamilies(ctx, &ecs.ListTaskDefinitionFamiliesInput{
		Status: ecstypes.TaskDefinitionFamilyStatusActive,
	})
	if err != nil {
		return nil, fmt.Errorf("list task def families: %w", err)
	}

	var defs []TaskDefinition
	for _, family := range listOut.Families {
		descOut, err := c.client.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
			TaskDefinition: &family,
		})
		if err != nil {
			continue
		}
		td := descOut.TaskDefinition
		var compat string
		for _, c := range td.Compatibilities {
			if compat != "" {
				compat += "/"
			}
			compat += string(c)
		}
		defs = append(defs, TaskDefinition{
			Family:        aws.ToString(td.Family),
			Revision:      td.Revision,
			ARN:           aws.ToString(td.TaskDefinitionArn),
			Status:        string(td.Status),
			CPU:           aws.ToString(td.Cpu),
			Memory:        aws.ToString(td.Memory),
			Compatibility: compat,
			RegisteredAt:  aws.ToTime(td.RegisteredAt),
		})
	}
	return defs, nil
}

// --- Service Events ---

type ServiceEvent struct {
	ID        string
	CreatedAt time.Time
	Message   string
}

func (c *ECSClient) GetServiceEvents(ctx context.Context, clusterARN, serviceName string) ([]ServiceEvent, error) {
	descOut, err := c.client.DescribeServices(ctx, &ecs.DescribeServicesInput{
		Cluster:  &clusterARN,
		Services: []string{serviceName},
	})
	if err != nil {
		return nil, fmt.Errorf("describe service events: %w", err)
	}
	if len(descOut.Services) == 0 {
		return nil, nil
	}

	var events []ServiceEvent
	for _, e := range descOut.Services[0].Events {
		events = append(events, ServiceEvent{
			ID:        aws.ToString(e.Id),
			CreatedAt: aws.ToTime(e.CreatedAt),
			Message:   aws.ToString(e.Message),
		})
	}
	return events, nil
}

// --- Stopped Tasks ---

type StoppedTask struct {
	TaskARN       string
	TaskID        string
	Status        string
	StoppedReason string
	StoppedAt     time.Time
	TaskDef       string
	Group         string
}

func (c *ECSClient) ListStoppedTasks(ctx context.Context, clusterARN string) ([]StoppedTask, error) {
	listOut, err := c.client.ListTasks(ctx, &ecs.ListTasksInput{
		Cluster:       &clusterARN,
		DesiredStatus: ecstypes.DesiredStatusStopped,
	})
	if err != nil {
		return nil, fmt.Errorf("list stopped tasks: %w", err)
	}
	if len(listOut.TaskArns) == 0 {
		return nil, nil
	}

	descOut, err := c.client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: &clusterARN,
		Tasks:   listOut.TaskArns,
	})
	if err != nil {
		return nil, fmt.Errorf("describe stopped tasks: %w", err)
	}

	var tasks []StoppedTask
	for _, t := range descOut.Tasks {
		arn := aws.ToString(t.TaskArn)
		tasks = append(tasks, StoppedTask{
			TaskARN:       arn,
			TaskID:        shortTaskID(arn),
			Status:        aws.ToString(t.LastStatus),
			StoppedReason: aws.ToString(t.StoppedReason),
			StoppedAt:     aws.ToTime(t.StoppedAt),
			TaskDef:       shortTaskDef(aws.ToString(t.TaskDefinitionArn)),
			Group:         aws.ToString(t.Group),
		})
	}
	return tasks, nil
}

// --- Service Deployments (for Resource Map) ---

type Deployment struct {
	ID             string
	Status         string
	TaskDefinition string
	DesiredCount   int32
	RunningCount   int32
	RolloutState   string
	CreatedAt      time.Time
}

func (c *ECSClient) GetServiceDeployments(ctx context.Context, clusterARN, serviceName string) ([]Deployment, []string, error) {
	descOut, err := c.client.DescribeServices(ctx, &ecs.DescribeServicesInput{
		Cluster:  &clusterARN,
		Services: []string{serviceName},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("describe service: %w", err)
	}
	if len(descOut.Services) == 0 {
		return nil, nil, nil
	}

	svc := descOut.Services[0]
	var deps []Deployment
	for _, d := range svc.Deployments {
		rollout := ""
		if d.RolloutState != "" {
			rollout = string(d.RolloutState)
		}
		deps = append(deps, Deployment{
			ID:             aws.ToString(d.Id),
			Status:         aws.ToString(d.Status),
			TaskDefinition: shortTaskDef(aws.ToString(d.TaskDefinition)),
			DesiredCount:   d.DesiredCount,
			RunningCount:   d.RunningCount,
			RolloutState:   rollout,
			CreatedAt:      aws.ToTime(d.CreatedAt),
		})
	}

	var tgARNs []string
	for _, lb := range svc.LoadBalancers {
		tgARNs = append(tgARNs, aws.ToString(lb.TargetGroupArn))
	}
	return deps, tgARNs, nil
}

// --- Mutation ---

func (c *ECSClient) UpdateServiceScale(ctx context.Context, clusterARN, serviceName string, desired int32) error {
	_, err := c.client.UpdateService(ctx, &ecs.UpdateServiceInput{
		Cluster:      &clusterARN,
		Service:      &serviceName,
		DesiredCount: &desired,
	})
	return err
}

func (c *ECSClient) ForceNewDeployment(ctx context.Context, clusterARN, serviceName string) error {
	_, err := c.client.UpdateService(ctx, &ecs.UpdateServiceInput{
		Cluster:            &clusterARN,
		Service:            &serviceName,
		ForceNewDeployment: true,
	})
	return err
}

func (c *ECSClient) StopTask(ctx context.Context, clusterARN, taskARN string) error {
	reason := "Stopped by ecs9s"
	_, err := c.client.StopTask(ctx, &ecs.StopTaskInput{
		Cluster: &clusterARN,
		Task:    &taskARN,
		Reason:  &reason,
	})
	return err
}

func (c *ECSClient) DeregisterTaskDefinition(ctx context.Context, arn string) error {
	_, err := c.client.DeregisterTaskDefinition(ctx, &ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: &arn,
	})
	return err
}

// --- ECS Exec ---

// ExecCheck verifies whether a task has ECS Exec enabled and returns diagnostics.
type ExecCheckResult struct {
	ClusterEnabled bool
	AgentRunning   bool   // managed agent RUNNING status
	TaskExecEnable bool   // task-level enableExecuteCommand
	InitStatus     string // managed agent lastStatus
	Details        string
}

func (c *ECSClient) CheckExecEnabled(ctx context.Context, clusterARN, taskARN string) (*ExecCheckResult, error) {
	descOut, err := c.client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: &clusterARN,
		Tasks:   []string{taskARN},
	})
	if err != nil {
		return nil, fmt.Errorf("describe task for exec check: %w", err)
	}
	if len(descOut.Tasks) == 0 {
		return nil, fmt.Errorf("task not found: %s", taskARN)
	}

	task := descOut.Tasks[0]
	result := &ExecCheckResult{
		TaskExecEnable: task.EnableExecuteCommand,
	}

	// Check the ExecuteCommandAgent managed agent status
	for _, ctr := range task.Containers {
		for _, ma := range ctr.ManagedAgents {
			if ma.Name == ecstypes.ManagedAgentNameExecuteCommandAgent {
				result.InitStatus = aws.ToString(ma.LastStatus)
				result.AgentRunning = result.InitStatus == "RUNNING"
			}
		}
	}

	if !result.TaskExecEnable {
		result.Details = "enableExecuteCommand is not set on this task. Update the service with --enable-execute-command."
	} else if !result.AgentRunning {
		result.Details = fmt.Sprintf("ExecuteCommand agent status: %s. The task may need to be restarted after enabling.", result.InitStatus)
	} else {
		result.Details = "ECS Exec is ready."
	}

	return result, nil
}

// EnableExecOnService enables ECS Exec on a service and triggers a new deployment
// so that new tasks start with the ExecuteCommandAgent.
func (c *ECSClient) EnableExecOnService(ctx context.Context, clusterARN, serviceName string) error {
	enable := true
	_, err := c.client.UpdateService(ctx, &ecs.UpdateServiceInput{
		Cluster:              &clusterARN,
		Service:              &serviceName,
		EnableExecuteCommand: &enable,
		ForceNewDeployment:   true,
	})
	if err != nil {
		return fmt.Errorf("enable execute command: %w", err)
	}
	return nil
}

// --- Task Definition Resources (for cost estimation) ---

type TaskDefResources struct {
	Family   string
	Revision int32
	CPU      string // e.g. "256", "512", "1024"
	Memory   string // e.g. "512", "1024", "2048"
}

func (c *ECSClient) GetTaskDefResources(ctx context.Context, taskDef string) (*TaskDefResources, error) {
	out, err := c.client.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &taskDef,
	})
	if err != nil {
		return nil, fmt.Errorf("describe task definition %s: %w", taskDef, err)
	}
	td := out.TaskDefinition
	return &TaskDefResources{
		Family:   aws.ToString(td.Family),
		Revision: td.Revision,
		CPU:      aws.ToString(td.Cpu),
		Memory:   aws.ToString(td.Memory),
	}, nil
}

// --- Helpers ---

func shortTaskDef(arn string) string {
	for i := len(arn) - 1; i >= 0; i-- {
		if arn[i] == '/' {
			return arn[i+1:]
		}
	}
	return arn
}

func shortTaskID(arn string) string {
	for i := len(arn) - 1; i >= 0; i-- {
		if arn[i] == '/' {
			return arn[i+1:]
		}
	}
	return arn
}
