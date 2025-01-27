package vaulty

import (
	"context"
	"fmt"
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

// WithGeneratedVaultClient creates a vault client with the given address.
//
// Deprecated: Use WithAddr instead for the same effect.
func WithGeneratedVaultClient(vaultAddress string) ClientOption {
	return WithAddr(vaultAddress)
}

func WithAddr(addr string) ClientOption {
	return func(c *client) {
		c.config.Address = addr
	}
}

func WithConfig(config *hashiVault.Config) ClientOption {
	return func(c *client) {
		c.config = config
	}
}

func WithTokenAuth(token string) ClientOption {
	return func(c *client) {
		c.auth = func(v *hashiVault.Client) (*hashiVault.Secret, error) {
			return tokenLogin(v, token)
		}
	}
}

func WithAppRoleAuth(roleID, secretID string) ClientOption {
	return func(c *client) {
		c.auth = func(v *hashiVault.Client) (*hashiVault.Secret, error) {
			sec, err := appRoleLogin(v, roleID, secretID)
			if err != nil {
				return nil, err
			}
			go c.renewAuthInfo()
			return sec, nil
		}
	}
}

func WithUserPassAuth(username, password string) ClientOption {
	return func(c *client) {
		c.auth = func(v *hashiVault.Client) (*hashiVault.Secret, error) {
			sec, err := userPassLogin(v, username, password)
			if err != nil {
				return nil, err
			}
			go c.renewAuthInfo()
			return sec, nil
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
