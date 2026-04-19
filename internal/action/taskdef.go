package action

import (
	"context"
	"fmt"

	"github.com/whchoi98/ecs9s/internal/aws"
)

func DeregisterTaskDefinition(ecs *aws.ECSClient, arn string) error {
	err := ecs.DeregisterTaskDefinition(context.Background(), arn)
	if err != nil {
		return fmt.Errorf("deregister task definition failed: %w", err)
	}
	return nil
}
