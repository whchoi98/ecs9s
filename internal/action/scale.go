package action

import (
	"context"
	"fmt"

	"github.com/whchoi98/ecs9s/internal/aws"
)

func ScaleService(ecs *aws.ECSClient, clusterARN, serviceName string, desired int32) error {
	err := ecs.UpdateServiceScale(context.Background(), clusterARN, serviceName, desired)
	if err != nil {
		return fmt.Errorf("scale service failed: %w", err)
	}
	return nil
}
