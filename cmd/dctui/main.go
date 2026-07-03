package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/openqt/tcmd/internal/app"
	"github.com/openqt/tcmd/internal/config"
)

func main() {
	if _, err := config.EnsureDir(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: config dir: %v\n", err)
	}

	leftPath, rightPath := ".", "."
	if wd, err := os.Getwd(); err == nil {
		leftPath = wd
		rightPath = wd
	}
	if len(os.Args) > 1 {
		leftPath = os.Args[1]
	}
	if len(os.Args) > 2 {
		rightPath = os.Args[2]
	}

	model := app.NewModel(leftPath, rightPath)
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
