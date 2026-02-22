package flow

import (
	"log/slog"

	"github.com/makehlv/code/clients"
	"github.com/makehlv/code/config"
)

type CodeFlowManageService struct {
	logger  *slog.Logger
	clients *clients.Clients
	config  *config.Config
}

func NewCodeFlowManageService(
	clients *clients.Clients, logger *slog.Logger, config *config.Config) *CodeFlowManageService {
	return &CodeFlowManageService{clients: clients, logger: logger, config: config}
}
