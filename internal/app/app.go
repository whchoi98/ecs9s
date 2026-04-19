package app

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/whchoi98/ecs9s/internal/action"
	"github.com/whchoi98/ecs9s/internal/aws"
	"github.com/whchoi98/ecs9s/internal/config"
	"github.com/whchoi98/ecs9s/internal/theme"
	"github.com/whchoi98/ecs9s/internal/ui"
	"github.com/whchoi98/ecs9s/internal/ui/components"
	"github.com/whchoi98/ecs9s/internal/ui/pages"
)

// execReadyMsg is sent after ECS Exec prerequisites pass.
type execReadyMsg struct {
	ClusterARN string
	TaskARN    string
	Container  string
}

type App struct {
	// UI components
	tabs       components.Tabs
	statusBar  components.StatusBar
	commandBar components.CommandBar
	help       components.Help
	confirm    components.Confirm

	// Pages
	clusterPage   pages.ClusterPage
	servicePage   pages.ServicePage
	taskPage      pages.TaskPage
	containerPage pages.ContainerPage
	taskDefPage   pages.TaskDefPage
	logsPage      pages.LogsPage
	ecrPage       pages.ECRPage
	elbPage       pages.ELBPage
	asgPage       pages.ASGPage
	vpcPage       pages.VPCPage
	iamPage       pages.IAMPage
	metricsPage   pages.MetricsPage
	ec2Page       pages.EC2Page
	eventsPage    pages.EventsPage
	stoppedPage   pages.StoppedPage
	resmapPage    pages.ResMapPage
	costPage      pages.CostPage
	ssmPage       pages.SSMPage
	secretsPage   pages.SecretsPage
	deployPage    pages.DeployPage
	alarmsPage    pages.AlarmsPage

	// State
	activePage ui.PageType
	navStack   []navEntry
	nav        ui.NavContext
	cfg        *config.Config
	session    *aws.Session
	thm        theme.Theme

	width, height int
}

type navEntry struct {
	page ui.PageType
	nav  ui.NavContext
}

func New(cfg *config.Config, session *aws.Session) App {
	thm := theme.Get(cfg.Theme)

	ecsClient := aws.NewECSClient(session.Config)
	cwClient := aws.NewCloudWatchClient(session.Config)
	ecrClient := aws.NewECRClient(session.Config)
	elbClient := aws.NewELBClient(session.Config)
	ec2Client := aws.NewEC2Client(session.Config)
	iamClient := aws.NewIAMClient(session.Config)
	asgClient := aws.NewAutoScalingClient(session.Config)
	ssmClient := aws.NewSSMClient(session.Config)
	secretsClient := aws.NewSecretsClient(session.Config)

	tabItems := []components.TabItem{
		{Name: "Cluster", Command: "cluster"},
		{Name: "Service", Command: "service"},
		{Name: "Task", Command: "task"},
		{Name: "Container", Command: "container"},
		{Name: "TaskDef", Command: "taskdef"},
		{Name: "Logs", Command: "log"},
		{Name: "ECR", Command: "ecr"},
		{Name: "ELB", Command: "elb"},
		{Name: "ASG", Command: "asg"},
		{Name: "VPC", Command: "vpc"},
		{Name: "IAM", Command: "iam"},
		{Name: "Metrics", Command: "metrics"},
		{Name: "EC2", Command: "ec2"},
		{Name: "Events", Command: "events"},
		{Name: "Stopped", Command: "stopped"},
		{Name: "ResMap", Command: "resmap"},
		{Name: "Cost", Command: "cost"},
		{Name: "SSM", Command: "ssm"},
		{Name: "Secrets", Command: "secrets"},
		{Name: "Deploy", Command: "deploy"},
		{Name: "Alarms", Command: "alarms"},
	}

	app := App{
		tabs:          components.NewTabs(tabItems, thm),
		statusBar:     components.NewStatusBar(thm),
		commandBar:    components.NewCommandBar(thm),
		help:          components.NewHelp(thm),
		confirm:       components.NewConfirm(thm),
		clusterPage:   pages.NewClusterPage(ecsClient, thm),
		servicePage:   pages.NewServicePage(ecsClient, thm),
		taskPage:      pages.NewTaskPage(ecsClient, thm),
		containerPage: pages.NewContainerPage(ecsClient, thm),
		taskDefPage:   pages.NewTaskDefPage(ecsClient, thm),
		logsPage:      pages.NewLogsPage(cwClient, thm),
		ecrPage:       pages.NewECRPage(ecrClient, thm),
		elbPage:       pages.NewELBPage(elbClient, thm),
		asgPage:       pages.NewASGPage(asgClient, thm),
		vpcPage:       pages.NewVPCPage(ec2Client, thm),
		iamPage:       pages.NewIAMPage(iamClient, thm),
		metricsPage:   pages.NewMetricsPage(cwClient, thm),
		ec2Page:       pages.NewEC2Page(ec2Client, thm),
		eventsPage:    pages.NewEventsPage(ecsClient, thm),
		stoppedPage:   pages.NewStoppedPage(ecsClient, thm),
		resmapPage:    pages.NewResMapPage(ecsClient, elbClient, thm),
		costPage:      pages.NewCostPage(ecsClient, thm),
		ssmPage:       pages.NewSSMPage(ssmClient, thm),
		secretsPage:   pages.NewSecretsPage(secretsClient, thm),
		deployPage:    pages.NewDeployPage(ecsClient, thm),
		alarmsPage:    pages.NewAlarmsPage(cwClient, thm),
		activePage:    ui.PageCluster,
		cfg:           cfg,
		session:       session,
		thm:           thm,
	}

	// Set status bar fields here (not in Init) because Init uses a value receiver
	app.statusBar.Profile = session.Profile
	app.statusBar.Region = session.Region
	return app
}

func (a App) Init() tea.Cmd {
	return a.clusterPage.Init()
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle confirm dialog first
	if a.confirm.Visible {
		var cmd tea.Cmd
		a.confirm, cmd = a.confirm.Update(msg)
		return a, cmd
	}

	// Handle command bar
	if a.commandBar.Active() {
		var cmd tea.Cmd
		a.commandBar, cmd = a.commandBar.Update(msg)
		return a, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.resize()
		return a, nil

	case tea.KeyMsg:
		if a.help.Visible {
			if msg.String() == "?" || msg.String() == "esc" {
				a.help.Toggle()
			}
			return a, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case "?":
			a.help.Toggle()
			return a, nil
		case ":":
			a.commandBar.Activate("command")
			return a, nil
		case "/":
			// Forward to active page
		case "tab", "]":
			a.tabs.Next()
			return a, a.switchPage(ui.PageType(a.tabs.Active))
		case "shift+tab", "[":
			a.tabs.Prev()
			return a, a.switchPage(ui.PageType(a.tabs.Active))
		case "R":
			return a, a.refreshCurrentPage()
		case "esc", "backspace":
			if len(a.navStack) > 0 {
				return a, a.goBack()
			}
		}

	case components.CommandExecuteMsg:
		if pt, ok := ui.PageTypeFromCommand(msg.Command); ok {
			a.tabs.SetActive(int(pt))
			return a, a.switchPage(pt)
		}
		a.statusBar.SetError(fmt.Sprintf("Unknown command: %s", msg.Command))
		return a, nil

	case ui.DrillDownMsg:
		a.navStack = append(a.navStack, navEntry{page: a.activePage, nav: a.nav})
		a.nav = msg.Context
		a.activePage = msg.Page
		a.tabs.SetActive(int(msg.Page))
		a.statusBar.Cluster = msg.Context.ClusterName
		return a, a.initCurrentPage()

	case ui.GoBackMsg:
		return a, a.goBack()

	case ui.ErrorMsg:
		a.statusBar.SetError(msg.Err.Error())
		return a, nil

	case ui.InfoMsg:
		a.statusBar.SetInfo(msg.Text)
		return a, nil

	case pages.ExecRequestMsg:
		// Step 1: Check local prerequisites (aws cli, session-manager-plugin)
		if err := action.CheckPrerequisites(); err != nil {
			a.statusBar.SetError(err.Error())
			return a, nil
		}

		// Step 2: Check if ECS Exec is enabled on the task
		ecsClient := aws.NewECSClient(a.session.Config)
		clusterARN := msg.ClusterARN
		taskARN := msg.TaskARN
		container := msg.Container
		return a, func() tea.Msg {
			check, err := ecsClient.CheckExecEnabled(context.Background(), clusterARN, taskARN)
			if err != nil {
				return ui.ErrorMsg{Err: fmt.Errorf("exec check: %w", err)}
			}
			if !check.TaskExecEnable || !check.AgentRunning {
				return ui.ErrorMsg{Err: fmt.Errorf("%s", check.Details)}
			}
			// All checks passed — return a message to start the shell
			return execReadyMsg{ClusterARN: clusterARN, TaskARN: taskARN, Container: container}
		}

	case execReadyMsg:
		// Step 3: Suspend TUI and run interactive shell
		cmd := action.ExecCommand(msg.ClusterARN, msg.TaskARN, msg.Container, "/bin/sh")
		return a, tea.ExecProcess(cmd, func(err error) tea.Msg {
			if err != nil {
				return ui.ErrorMsg{Err: fmt.Errorf("ECS Exec: %w", err)}
			}
			return ui.InfoMsg{Text: fmt.Sprintf("Exited shell: %s", msg.Container)}
		})
	}

	// Forward to active page
	cmd := a.updateActivePage(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a App) View() string {
	if a.width == 0 {
		return "  Loading ecs9s..."
	}

	tabBar := a.tabs.View()
	statusBar := a.statusBar.View()

	pageHeight := a.height - lipgloss.Height(tabBar) - lipgloss.Height(statusBar) - 2
	if pageHeight < 1 {
		pageHeight = 1
	}

	pageContent := a.viewActivePage()
	page := lipgloss.NewStyle().
		Height(pageHeight).
		Width(a.width).
		Render(pageContent)

	cmdBar := a.commandBar.View()

	screen := lipgloss.JoinVertical(lipgloss.Left,
		tabBar,
		statusBar,
		page,
		cmdBar,
	)

	// Overlay help or confirm
	if a.help.Visible {
		return a.help.View()
	}
	if a.confirm.Visible {
		return a.confirm.View()
	}

	return screen
}

// --- Internal ---

func (a *App) resize() {
	a.tabs.SetWidth(a.width)
	a.statusBar.SetWidth(a.width)
	a.commandBar.SetWidth(a.width)
	a.help.SetSize(a.width, a.height)
	a.confirm.SetSize(a.width, a.height)

	pageH := a.height - 6
	if pageH < 1 {
		pageH = 1
	}
	a.clusterPage.SetSize(a.width, pageH)
	a.servicePage.SetSize(a.width, pageH)
	a.taskPage.SetSize(a.width, pageH)
	a.containerPage.SetSize(a.width, pageH)
	a.taskDefPage.SetSize(a.width, pageH)
	a.logsPage.SetSize(a.width, pageH)
	a.ecrPage.SetSize(a.width, pageH)
	a.elbPage.SetSize(a.width, pageH)
	a.asgPage.SetSize(a.width, pageH)
	a.vpcPage.SetSize(a.width, pageH)
	a.iamPage.SetSize(a.width, pageH)
	a.metricsPage.SetSize(a.width, pageH)
	a.ec2Page.SetSize(a.width, pageH)
	a.eventsPage.SetSize(a.width, pageH)
	a.stoppedPage.SetSize(a.width, pageH)
	a.resmapPage.SetSize(a.width, pageH)
	a.costPage.SetSize(a.width, pageH)
	a.ssmPage.SetSize(a.width, pageH)
	a.secretsPage.SetSize(a.width, pageH)
	a.deployPage.SetSize(a.width, pageH)
	a.alarmsPage.SetSize(a.width, pageH)
}

func (a *App) resetContext() {
	a.nav = ui.NavContext{}
	a.navStack = nil
	a.statusBar.Cluster = ""
}

func (a *App) switchPage(pt ui.PageType) tea.Cmd {
	a.activePage = pt
	a.statusBar.ClearMessage()

	// Global pages: hide cluster in status bar (they show all resources)
	// Context-dependent pages: show current drill-down cluster if available
	switch pt {
	case ui.PageCluster, ui.PageTaskDef, ui.PageECR, ui.PageELB,
		ui.PageVPC, ui.PageIAM, ui.PageEC2, ui.PageSSM,
		ui.PageSecrets, ui.PageAlarms, ui.PageASG:
		a.statusBar.Cluster = ""
	default:
		a.statusBar.Cluster = a.nav.ClusterName
	}

	return a.initCurrentPage()
}

func (a *App) initCurrentPage() tea.Cmd {
	switch a.activePage {
	case ui.PageCluster:
		return a.clusterPage.Refresh()
	case ui.PageService:
		return a.servicePage.SetContext(a.nav)
	case ui.PageTask:
		return a.taskPage.SetContext(a.nav)
	case ui.PageContainer:
		return a.containerPage.SetContext(a.nav)
	case ui.PageTaskDef:
		return a.taskDefPage.Refresh()
	case ui.PageLogs:
		return nil // needs log group selection
	case ui.PageECR:
		return a.ecrPage.Refresh()
	case ui.PageELB:
		return a.elbPage.Refresh()
	case ui.PageASG:
		return a.asgPage.Refresh()
	case ui.PageVPC:
		return a.vpcPage.Refresh()
	case ui.PageIAM:
		return a.iamPage.Refresh()
	case ui.PageMetrics:
		if a.nav.ServiceName != "" {
			return a.metricsPage.SetContext(a.nav)
		}
		return nil
	case ui.PageEC2:
		return a.ec2Page.Refresh()
	case ui.PageEvents:
		if a.nav.ServiceName != "" {
			return a.eventsPage.SetContext(a.nav)
		}
		return nil
	case ui.PageStopped:
		if a.nav.ClusterARN != "" {
			return a.stoppedPage.SetContext(a.nav)
		}
		return nil
	case ui.PageResMap:
		if a.nav.ServiceName != "" {
			return a.resmapPage.SetContext(a.nav)
		}
		return nil
	case ui.PageCost:
		if a.nav.ClusterARN != "" {
			return a.costPage.SetContext(a.nav)
		}
		return nil
	case ui.PageSSM:
		return a.ssmPage.Refresh()
	case ui.PageSecrets:
		return a.secretsPage.Refresh()
	case ui.PageDeploy:
		if a.nav.ServiceName != "" {
			return a.deployPage.SetContext(a.nav)
		}
		return nil
	case ui.PageAlarms:
		return a.alarmsPage.Refresh()
	}
	return nil
}

func (a *App) refreshCurrentPage() tea.Cmd {
	return a.initCurrentPage()
}

func (a *App) goBack() tea.Cmd {
	if len(a.navStack) == 0 {
		return nil
	}
	prev := a.navStack[len(a.navStack)-1]
	a.navStack = a.navStack[:len(a.navStack)-1]
	a.activePage = prev.page
	a.nav = prev.nav
	a.tabs.SetActive(int(prev.page))
	if prev.nav.ClusterName != "" {
		a.statusBar.Cluster = prev.nav.ClusterName
	} else {
		a.statusBar.Cluster = ""
	}
	return a.initCurrentPage()
}

func (a *App) updateActivePage(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch a.activePage {
	case ui.PageCluster:
		a.clusterPage, cmd = a.clusterPage.Update(msg)
	case ui.PageService:
		a.servicePage, cmd = a.servicePage.Update(msg)
	case ui.PageTask:
		a.taskPage, cmd = a.taskPage.Update(msg)
	case ui.PageContainer:
		a.containerPage, cmd = a.containerPage.Update(msg)
	case ui.PageTaskDef:
		a.taskDefPage, cmd = a.taskDefPage.Update(msg)
	case ui.PageLogs:
		a.logsPage, cmd = a.logsPage.Update(msg)
	case ui.PageECR:
		a.ecrPage, cmd = a.ecrPage.Update(msg)
	case ui.PageELB:
		a.elbPage, cmd = a.elbPage.Update(msg)
	case ui.PageASG:
		a.asgPage, cmd = a.asgPage.Update(msg)
	case ui.PageVPC:
		a.vpcPage, cmd = a.vpcPage.Update(msg)
	case ui.PageIAM:
		a.iamPage, cmd = a.iamPage.Update(msg)
	case ui.PageMetrics:
		a.metricsPage, cmd = a.metricsPage.Update(msg)
	case ui.PageEC2:
		a.ec2Page, cmd = a.ec2Page.Update(msg)
	case ui.PageEvents:
		a.eventsPage, cmd = a.eventsPage.Update(msg)
	case ui.PageStopped:
		a.stoppedPage, cmd = a.stoppedPage.Update(msg)
	case ui.PageResMap:
		a.resmapPage, cmd = a.resmapPage.Update(msg)
	case ui.PageCost:
		a.costPage, cmd = a.costPage.Update(msg)
	case ui.PageSSM:
		a.ssmPage, cmd = a.ssmPage.Update(msg)
	case ui.PageSecrets:
		a.secretsPage, cmd = a.secretsPage.Update(msg)
	case ui.PageDeploy:
		a.deployPage, cmd = a.deployPage.Update(msg)
	case ui.PageAlarms:
		a.alarmsPage, cmd = a.alarmsPage.Update(msg)
	}
	return cmd
}

func (a App) viewActivePage() string {
	switch a.activePage {
	case ui.PageCluster:
		return a.clusterPage.View()
	case ui.PageService:
		return a.servicePage.View()
	case ui.PageTask:
		return a.taskPage.View()
	case ui.PageContainer:
		return a.containerPage.View()
	case ui.PageTaskDef:
		return a.taskDefPage.View()
	case ui.PageLogs:
		return a.logsPage.View()
	case ui.PageECR:
		return a.ecrPage.View()
	case ui.PageELB:
		return a.elbPage.View()
	case ui.PageASG:
		return a.asgPage.View()
	case ui.PageVPC:
		return a.vpcPage.View()
	case ui.PageIAM:
		return a.iamPage.View()
	case ui.PageMetrics:
		return a.metricsPage.View()
	case ui.PageEC2:
		return a.ec2Page.View()
	case ui.PageEvents:
		return a.eventsPage.View()
	case ui.PageStopped:
		return a.stoppedPage.View()
	case ui.PageResMap:
		return a.resmapPage.View()
	case ui.PageCost:
		return a.costPage.View()
	case ui.PageSSM:
		return a.ssmPage.View()
	case ui.PageSecrets:
		return a.secretsPage.View()
	case ui.PageDeploy:
		return a.deployPage.View()
	case ui.PageAlarms:
		return a.alarmsPage.View()
	}
	return ""
}
