package vaulty

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	hashiVault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"
)

type ClientOption func(c *client)

func WithContext(ctx context.Context) ClientOption {
	return func(c *client) {
		c.ctx = ctx
	}
}

func WithGeneratedVaultClient(vaultAddress string) ClientOption {
	return func(c *client) {
		config := hashiVault.DefaultConfig()
		config.Address = vaultAddress

		vc, err := hashiVault.NewClient(config)
		if err != nil {
			slog.Error("Error creating vault client", slog.String(loggingKeyError, err.Error()))
			os.Exit(1)
		}

		c.v = vc
	}
}

func WithAppRoleAuth(roleID, secretID string) ClientOption {
	return func(c *client) {
		c.auth = func(v *hashiVault.Client) (*hashiVault.Secret, error) {
			return appRoleLogin(v, roleID, secretID)
		}
	}
}

func WithUserPassAuth(username, password string) ClientOption {
	return func(c *client) {
		c.auth = func(v *hashiVault.Client) (*hashiVault.Secret, error) {
			return userPassLogin(v, username, password)
		}
	}
}

func WithKvv2Mount(mount string) ClientOption {
	return func(c *client) {
		c.kvv2Mount = mount
	}
}

func WithKubernetesAuthDefault() ClientOption {
	return func(c *client) {
		c.auth = func(v *hashiVault.Client) (*hashiVault.Secret, error) {
			role := os.Getenv(envKubernetesRole)
			if role == "" {
				return nil, fmt.Errorf("%s environment variable not set", envKubernetesRole)
			}

			return kubernetesLogin(v, role, auth.WithServiceAccountTokenPath(kubernetesServiceAccountTokenPath))
		}
	}
}

func WithKubernetesAuthFromEnv() ClientOption {
	return func(c *client) {
		c.auth = func(v *hashiVault.Client) (*hashiVault.Secret, error) {
			role := os.Getenv(envKubernetesRole)
			if role == "" {
				return nil, fmt.Errorf("%s environment variable not set", envKubernetesRole)
			}

			return kubernetesLogin(v, role, auth.WithServiceAccountTokenEnv(envKubernetesToken))
		}
	}
}

func WithKubernetesAuth(role, token string) ClientOption {
	return func(c *client) {
		c.auth = func(v *hashiVault.Client) (*hashiVault.Secret, error) {
			return kubernetesLogin(v, role, auth.WithServiceAccountToken(token))
		}
	}
}
