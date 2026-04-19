package ui

import "testing"

func TestPageTypeFromCommand(t *testing.T) {
	tests := []struct {
		cmd  string
		want PageType
		ok   bool
	}{
		{"cluster", PageCluster, true},
		{"service", PageService, true},
		{"task", PageTask, true},
		{"container", PageContainer, true},
		{"taskdef", PageTaskDef, true},
		{"log", PageLogs, true},
		{"ecr", PageECR, true},
		{"elb", PageELB, true},
		{"asg", PageASG, true},
		{"vpc", PageVPC, true},
		{"iam", PageIAM, true},
		{"metrics", PageMetrics, true},
		{"ec2", PageEC2, true},
		{"events", PageEvents, true},
		{"stopped", PageStopped, true},
		{"resmap", PageResMap, true},
		{"cost", PageCost, true},
		{"ssm", PageSSM, true},
		{"secrets", PageSecrets, true},
		{"deploy", PageDeploy, true},
		{"alarms", PageAlarms, true},
		{"unknown", 0, false},
		{"", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.cmd, func(t *testing.T) {
			got, ok := PageTypeFromCommand(tt.cmd)
			if ok != tt.ok {
				t.Errorf("PageTypeFromCommand(%q) ok = %v, want %v", tt.cmd, ok, tt.ok)
			}
			if ok && got != tt.want {
				t.Errorf("PageTypeFromCommand(%q) = %v, want %v", tt.cmd, got, tt.want)
			}
		})
	}
}

func TestPageTypeString(t *testing.T) {
	if PageCluster.String() != "Cluster" {
		t.Errorf("PageCluster.String() = %q", PageCluster.String())
	}
	if PageAlarms.String() != "Alarms" {
		t.Errorf("PageAlarms.String() = %q", PageAlarms.String())
	}
}

func TestPageTypeCommandRoundtrip(t *testing.T) {
	for i := PageCluster; i <= PageAlarms; i++ {
		cmd := i.Command()
		got, ok := PageTypeFromCommand(cmd)
		if !ok {
			t.Errorf("PageTypeFromCommand(%q) returned false for valid PageType %d", cmd, i)
			continue
		}
		if got != i {
			t.Errorf("roundtrip: PageType %d → Command %q → PageType %d", i, cmd, got)
		}
	}
}
