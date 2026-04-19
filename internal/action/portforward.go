package action

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// PortForward starts an SSM port forwarding session to an ECS container.
// clusterName: short cluster name (not full ARN)
// taskID: short task ID (not full ARN)
// runtimeID: container runtime ID from ECS DescribeTasks
func PortForward(clusterName, taskID, runtimeID string, localPort, remotePort int) error {
	// SSM target format for ECS: ecs:{cluster-name}_{task-id}_{runtime-id}
	target := fmt.Sprintf("ecs:%s_%s_%s", clusterName, taskID, runtimeID)

	args := []string{
		"ssm", "start-session",
		"--target", target,
		"--document-name", "AWS-StartPortForwardingSession",
		"--parameters", fmt.Sprintf(`{"portNumber":["%d"],"localPortNumber":["%d"]}`, remotePort, localPort),
	}

	cmd := exec.Command("aws", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("port forward failed: %w", err)
	}
	return nil
}

// ExtractClusterName extracts the short cluster name from a full ARN.
// e.g. "arn:aws:ecs:us-east-1:123456:cluster/my-cluster" -> "my-cluster"
func ExtractClusterName(arn string) string {
	if idx := strings.LastIndex(arn, "/"); idx != -1 {
		return arn[idx+1:]
	}
	return arn
}

// ExtractTaskID extracts the short task ID from a full task ARN.
// e.g. "arn:aws:ecs:us-east-1:123456:task/my-cluster/abc123" -> "abc123"
func ExtractTaskID(arn string) string {
	if idx := strings.LastIndex(arn, "/"); idx != -1 {
		return arn[idx+1:]
	}
	return arn
}
