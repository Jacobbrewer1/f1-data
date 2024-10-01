package vault

import (
	"context"
	"log/slog"
	"os"

	"github.com/Jacobbrewer1/f1-data/pkg/logging"
	vault "github.com/hashicorp/vault/api"
)

type ClientOption func(c *client)

func WithContext(ctx context.Context) ClientOption {
	return func(c *client) {
		c.ctx = ctx
	}
}

func WithGeneratedVaultClient(vaultAddress string) ClientOption {
	return func(c *client) {
		config := vault.DefaultConfig()
		config.Address = vaultAddress

		vc, err := vault.NewClient(config)
		if err != nil {
			slog.Error("Error creating vault client", slog.String(logging.KeyError, err.Error()))
			os.Exit(1)
		}

		c.v = vc
	}
}

func WithTransitEncrypt(path string) ClientOption {
	return func(c *client) {
		c.transitPathEncrypt = path
	}
}

func WithTransitDecrypt(path string) ClientOption {
	return func(c *client) {
		c.transitPathDecrypt = path
	}
}

func WithAppRoleAuth(roleID, secretID string) ClientOption {
	return func(c *client) {
		c.auth = func(v *vault.Client) (*vault.Secret, error) {
			return appRoleLogin(v, roleID, secretID)
		}
	}
}

func WithUserPassAuth(username, password string) ClientOption {
	return func(c *client) {
		c.auth = func(v *vault.Client) (*vault.Secret, error) {
			return userPassLogin(v, username, password)
		}
	}
}

func WithKvv2Mount(mount string) ClientOption {
	return func(c *client) {
		c.kvv2Mount = mount
	}
}
