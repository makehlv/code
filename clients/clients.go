package clients

import (
	"github.com/makehlv/code/clients/git"
)

type Clients struct {
	Git *git.GitClient
}

func NewClients() *Clients {
	return &Clients{Git: git.NewGitClient()}
}
