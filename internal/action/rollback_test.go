package action

import "testing"

func TestPreviousRevision(t *testing.T) {
	tests := []struct {
		name    string
		taskDef string
		want    string
		wantErr bool
	}{
		{
			name:    "normal revision",
			taskDef: "my-app:5",
			want:    "my-app:4",
		},
		{
			name:    "full ARN with revision",
			taskDef: "arn:aws:ecs:us-east-1:123456:task-definition/my-app:10",
			want:    "arn:aws:ecs:us-east-1:123456:task-definition/my-app:9",
		},
		{
			name:    "revision 1 cannot go lower",
			taskDef: "my-app:1",
			wantErr: true,
		},
		{
			name:    "no colon separator",
			taskDef: "my-app",
			wantErr: true,
		},
		{
			name:    "non-numeric revision",
			taskDef: "my-app:latest",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := previousRevision(tt.taskDef)
			if tt.wantErr {
				if err == nil {
					t.Errorf("previousRevision(%q) expected error, got %q", tt.taskDef, got)
				}
				return
			}
			if err != nil {
				t.Errorf("previousRevision(%q) unexpected error: %v", tt.taskDef, err)
				return
			}
			if got != tt.want {
				t.Errorf("previousRevision(%q) = %q, want %q", tt.taskDef, got, tt.want)
			}
		})
	}
}
