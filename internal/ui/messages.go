package ui

// PageType identifies which page/view is active.
type PageType int

const (
	PageCluster PageType = iota
	PageService
	PageTask
	PageContainer
	PageTaskDef
	PageLogs
	PageECR
	PageELB
	PageASG
	PageVPC
	PageIAM
	PageMetrics
	PageEC2
	PageEvents
	PageStopped
	PageResMap
	PageCost
	PageSSM
	PageSecrets
	PageDeploy
	PageAlarms
)

func (p PageType) String() string {
	return [...]string{
		"Cluster", "Service", "Task", "Container", "TaskDef",
		"Logs", "ECR", "ELB", "ASG", "VPC", "IAM", "Metrics", "EC2",
		"Events", "Stopped", "ResMap", "Cost",
		"SSM", "Secrets", "Deploy", "Alarms",
	}[p]
}

func (p PageType) Command() string {
	return [...]string{
		"cluster", "service", "task", "container", "taskdef",
		"log", "ecr", "elb", "asg", "vpc", "iam", "metrics", "ec2",
		"events", "stopped", "resmap", "cost",
		"ssm", "secrets", "deploy", "alarms",
	}[p]
}

func PageTypeFromCommand(cmd string) (PageType, bool) {
	for i := PageCluster; i <= PageAlarms; i++ {
		if i.Command() == cmd {
			return i, true
		}
	}
	return 0, false
}

// NavContext holds the drill-down context as the user navigates deeper.
type NavContext struct {
	ClusterARN  string
	ClusterName string
	ServiceName string
	ServiceARN  string
	TaskARN     string
	TaskID      string
}

// --- Messages ---

type SwitchPageMsg struct {
	Page PageType
}

type DrillDownMsg struct {
	Page    PageType
	Context NavContext
}

type GoBackMsg struct{}

type RefreshMsg struct{}

type ErrorMsg struct {
	Err error
}

type InfoMsg struct {
	Text string
}

type DataLoadedMsg struct {
	Page PageType
	Data interface{}
}
