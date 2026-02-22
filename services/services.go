package services

import (
	"log/slog"

	"github.com/makehlv/code/clients"
	"github.com/makehlv/code/config"
	"github.com/makehlv/code/services/flow"
)

type Services struct {
	Flow *flow.CodeFlowManageService
}

func NewServices(clients *clients.Clients, logger *slog.Logger, config *config.Config) *Services {
	return &Services{Flow: flow.NewCodeFlowManageService(clients, logger, config)}
}
