package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/makehlv/code/clients"
	"github.com/makehlv/code/config"
	"github.com/makehlv/code/services"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: code <command> [flags]")
		os.Exit(1)
	}

	clients := clients.NewClients()
	logger := slog.New(NewColorHandler(os.Stderr, slog.LevelInfo))
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error("config failed", "error", err)
		os.Exit(1)
	}
	svc := services.NewServices(clients, logger, cfg)

	command := os.Args[1]
	switch command {
	case "squash":
		comparableBranch := parseFlag(os.Args[2:], "--compare")
		if comparableBranch == "" {
			comparableBranch = "develop"
		}
		message := parseFlag(os.Args[2:], "--message")
		if err := svc.Flow.Squash(comparableBranch, message, hasFlag(os.Args[2:], "--push-force")); err != nil {
			logger.Error("squash failed", "error", err)
			os.Exit(1)
		}
	case "clean":
		if err := svc.Flow.CleanFallbackBranches(); err != nil {
			logger.Error("clean failed", "error", err)
			os.Exit(1)
		}
	case "commit":
		if err := svc.Flow.Commit(); err != nil {
			logger.Error("commit failed", "error", err)
			os.Exit(1)
		}
	case "push":
		if err := svc.Flow.Push(); err != nil {
			logger.Error("push failed", "error", err)
			os.Exit(1)
		}
	case "jira":
		link, err := svc.Flow.JiraLink()
		if err != nil {
			logger.Error("jira failed", "error", err)
			os.Exit(1)
		}
		if hasFlag(os.Args[2:], "--link") {
			fmt.Println(link)
			return
		}
		svc.Flow.OpenInChromeOrPrint(link)
	case "gitlab":
		link, err := svc.Flow.GitlabLink(hasFlag(os.Args[2:], "--mr"))
		if err != nil {
			logger.Error("gitlab failed", "error", err)
			os.Exit(1)
		}
		if hasFlag(os.Args[2:], "--link") {
			fmt.Println(link)
			return
		}
		svc.Flow.OpenInChromeOrPrint(link)
	default:
		fmt.Printf("unknown command: %s\n", command)
		os.Exit(1)
	}
}

func parseFlag(args []string, flag string) string {
	for i, arg := range args {
		if arg == flag && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func hasFlag(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag {
			return true
		}
	}
	return false
}
