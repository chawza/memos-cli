package cmd

import (
	"github.com/chawza/memos-cli/internal/api"
)

// client wraps api.Client with a concrete type alias for convenience.
type client = api.Client

func newClient(baseURL, token string) *client {
	return (*client)(api.NewClient(baseURL, token))
}
