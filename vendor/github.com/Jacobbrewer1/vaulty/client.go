package vaulty

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	hashiVault "github.com/hashicorp/vault/api"
)

var (
	ErrSecretNotFound = hashiVault.ErrSecretNotFound
)

type ClientHandler interface {
	Client() *hashiVault.Client
}

type Client interface {
	ClientHandler

	// Path returns the secret path for the given name.
	Path(name string, opts ...PathOption) Repository
}

type (
	RenewalFunc func() (*hashiVault.Secret, error)
	loginFunc   func(v *hashiVault.Client) (*hashiVault.Secret, error)
)

type client struct {
	ctx context.Context

	transitPathEncrypt string
	transitPathDecrypt string

	kvv2Mount string

	auth loginFunc

	// Below are set on initialization
	v         *hashiVault.Client
	authCreds *hashiVault.Secret
}

func NewClient(opts ...ClientOption) (Client, error) {
	c := new(client)

	for _, opt := range opts {
		opt(c)
	}

	if c.ctx == nil {
		c.ctx = context.Background()
	}

	if c.v == nil {
		return nil, errors.New("vault client is nil")
	} else if c.auth == nil {
		return nil, errors.New("auth method is nil")
	}

	authCreds, err := c.auth(c.v)
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate with Vault: %w", err)
	}

	c.authCreds = authCreds

	go c.renewAuthInfo()

	return c, nil
}

func (c *client) renewAuthInfo() {
	err := RenewLease(c.ctx, c, "auth", c.authCreds, func() (*hashiVault.Secret, error) {
		authInfo, err := c.auth(c.v)
		if err != nil {
			return nil, fmt.Errorf("unable to renew auth info: %w", err)
		}

		c.authCreds = authInfo

		return authInfo, nil
	})
	if err != nil {
		slog.Error("unable to renew auth info", slog.String(loggingKeyError, err.Error()))
		os.Exit(1)
	}
}

func (c *client) Client() *hashiVault.Client {
	return c.v
}

func (c *client) Path(name string, opts ...PathOption) Repository {
	p := &SecretPath{
		r:     c,
		mount: c.kvv2Mount, // Default to kvv2
		name:  name,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}
