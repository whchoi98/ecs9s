package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/whchoi98/ecs9s/internal/app"
	"github.com/whchoi98/ecs9s/internal/aws"
	"github.com/whchoi98/ecs9s/internal/config"
)

var version = "dev"

func main() {
	profile := flag.String("profile", "", "AWS profile")
	region := flag.String("region", "", "AWS region")
	themeFlag := flag.String("theme", "", "Theme (dark, light, blue)")
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("ecs9s %s\n", version)
		os.Exit(0)
	}

	cfg := config.Load()

	// CLI flags override config
	if *profile != "" {
		cfg.AWS.Profile = *profile
	}
	if *region != "" {
		cfg.AWS.Region = *region
	}
	if *themeFlag != "" {
		cfg.Theme = *themeFlag
	}

	session, err := aws.NewSession(cfg.AWS.Profile, cfg.AWS.Region)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize AWS session: %v\n", err)
		os.Exit(1)
	}

	application := app.New(cfg, session)

	p := tea.NewProgram(
		application,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
