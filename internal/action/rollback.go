package action

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func RollbackService(cfg aws.Config, clusterARN, serviceName string) error {
	client := ecs.NewFromConfig(cfg)

	descOut, err := client.DescribeServices(context.Background(), &ecs.DescribeServicesInput{
		Cluster:  &clusterARN,
		Services: []string{serviceName},
	})
	if err != nil {
		return fmt.Errorf("describe service: %w", err)
	}
	if len(descOut.Services) == 0 {
		return fmt.Errorf("service not found: %s", serviceName)
	}

	svc := descOut.Services[0]

	// Find the PRIMARY deployment's task definition
	var primaryTaskDef string
	for _, d := range svc.Deployments {
		if strings.EqualFold(aws.ToString(d.Status), "PRIMARY") {
			primaryTaskDef = aws.ToString(d.TaskDefinition)
			break
		}
	}
	if primaryTaskDef == "" {
		return fmt.Errorf("no PRIMARY deployment found")
	}

	// Find the most recent non-PRIMARY deployment with a different task definition
	var rollbackTaskDef string
	for _, d := range svc.Deployments {
		td := aws.ToString(d.TaskDefinition)
		if td != primaryTaskDef {
			rollbackTaskDef = td
			break
		}
	}

	// If no different deployment exists, try the previous task definition revision
	if rollbackTaskDef == "" {
		prevRev, err := previousRevision(primaryTaskDef)
		if err != nil {
			return fmt.Errorf("no previous deployment to rollback to and cannot determine previous revision")
		}

		// Verify the previous revision exists
		_, descErr := client.DescribeTaskDefinition(context.Background(), &ecs.DescribeTaskDefinitionInput{
			TaskDefinition: &prevRev,
		})
		if descErr != nil {
			return fmt.Errorf("previous task definition revision %s not found: %w", prevRev, descErr)
		}
		rollbackTaskDef = prevRev
	}

	_, err = client.UpdateService(context.Background(), &ecs.UpdateServiceInput{
		Cluster:            &clusterARN,
		Service:            &serviceName,
		TaskDefinition:     &rollbackTaskDef,
		ForceNewDeployment: true,
	})
	if err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}
	return nil
}

// previousRevision decrements the revision number in a task definition ARN/family:revision string.
func previousRevision(taskDef string) (string, error) {
	idx := strings.LastIndex(taskDef, ":")
	if idx == -1 {
		return "", fmt.Errorf("invalid task definition format: %s", taskDef)
	}

	var rev int
	_, err := fmt.Sscanf(taskDef[idx+1:], "%d", &rev)
	if err != nil || rev <= 1 {
		return "", fmt.Errorf("cannot determine previous revision for: %s", taskDef)
	}

	return fmt.Sprintf("%s:%d", taskDef[:idx], rev-1), nil
}
