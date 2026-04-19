package action

import (
	"context"
	"fmt"

	"github.com/whchoi98/ecs9s/internal/aws"
)

func ForceNewDeployment(ecs *aws.ECSClient, clusterARN, serviceName string) error {
	err := ecs.ForceNewDeployment(context.Background(), clusterARN, serviceName)
	if err != nil {
		return fmt.Errorf("force deploy failed: %w", err)
	}
	return nil
}
