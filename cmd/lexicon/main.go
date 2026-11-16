package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kristyancarvalho/lexicon/internal/config"
	"github.com/kristyancarvalho/lexicon/internal/ui"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	if *showVersion {
		fmt.Println("lexicon", version)
		return
	}

	if err := run(*configPath); err != nil {
		fmt.Fprintln(os.Stderr, "lexicon:", err)
		os.Exit(1)
	}
}

func run(path string) error {
	if path == "" {
		resolved, err := config.DefaultPath()
		if err != nil {
			return err
		}
		path = resolved
	}

	cfg, err := config.Load(path)
	if err != nil {
		return err
	}

	_, err = tea.NewProgram(ui.New(cfg, path), tea.WithAltScreen()).Run()
	return err
}
