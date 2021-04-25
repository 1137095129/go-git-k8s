package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/wang1137095129/go-git-k8s/config"
)

type Handler interface {
	OpenRepository(c *config.Config) (*git.Repository, error)
}
