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
		fmt.Println("usage: kk <command> [flags]")
		os.Exit(1)
	}

	clients := clients.NewClients()
	logger := slog.New(NewColorHandler(os.Stderr, slog.LevelInfo))
	config := config.NewConfig()
	svc := services.NewServices(clients, logger, config)

	command := os.Args[1]
	switch command {
	case "squash":
		comparableBranch := parseFlag(os.Args[2:], "--compare")
		if comparableBranch == "" {
			comparableBranch = "develop"
		}
		message := parseFlag(os.Args[2:], "--message")
		if err := svc.Flow.Squash(comparableBranch, message); err != nil {
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
