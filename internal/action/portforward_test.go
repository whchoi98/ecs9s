package action

import "testing"

func TestExtractClusterName(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want string
	}{
		{"full ARN", "arn:aws:ecs:us-east-1:123456789:cluster/my-cluster", "my-cluster"},
		{"short name", "my-cluster", "my-cluster"},
		{"nested path", "arn:aws:ecs:ap-northeast-2:999:cluster/prod-ecs", "prod-ecs"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractClusterName(tt.arn)
			if got != tt.want {
				t.Errorf("ExtractClusterName(%q) = %q, want %q", tt.arn, got, tt.want)
			}
		})
	}
}

func TestExtractTaskID(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want string
	}{
		{"full ARN", "arn:aws:ecs:us-east-1:123456789:task/my-cluster/abc123def456", "abc123def456"},
		{"short ID", "abc123def456", "abc123def456"},
		{"only task path", "task/abc123", "abc123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractTaskID(tt.arn)
			if got != tt.want {
				t.Errorf("ExtractTaskID(%q) = %q, want %q", tt.arn, got, tt.want)
			}
		})
	}
}
