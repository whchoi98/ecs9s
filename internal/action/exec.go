package action

import (
	"fmt"
	"os/exec"
)

// Shells available for ECS Exec, in preference order.
var Shells = []string{"/bin/bash", "/bin/sh"}

// ExecCommand builds the *exec.Cmd for ECS Exec.
// The caller (app shell) uses tea.ExecProcess to run it,
// which suspends the TUI, hands stdin/stdout to the subprocess,
// and restores the TUI when the subprocess exits.
func ExecCommand(clusterARN, taskARN, containerName, shell string) *exec.Cmd {
	if shell == "" {
		shell = "/bin/sh"
	}

	return exec.Command("aws",
		"ecs", "execute-command",
		"--cluster", clusterARN,
		"--task", taskARN,
		"--container", containerName,
		"--interactive",
		"--command", shell,
	)
}

// CheckPrerequisites verifies that the local tools required for ECS Exec are installed.
func CheckPrerequisites() error {
	// Check AWS CLI
	if _, err := exec.LookPath("aws"); err != nil {
		return fmt.Errorf("AWS CLI not found in PATH. Install: https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html")
	}

	// Check session-manager-plugin (required for ECS Exec)
	if _, err := exec.LookPath("session-manager-plugin"); err != nil {
		return fmt.Errorf("session-manager-plugin not found. Install: https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html")
	}

	return nil
}
