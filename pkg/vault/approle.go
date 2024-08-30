package vault

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/Jacobbrewer1/f1-data/pkg/logging"
	vault "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/spf13/viper"
)

type appRoleClient struct {
	v        *vault.Client
	authInfo *vault.Secret
	vip      *viper.Viper
}

func NewClientAppRole(v *viper.Viper) (Client, error) {
	config := vault.DefaultConfig()
	config.Address = v.GetString("vault.address")

	c, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	clientImpl := &appRoleClient{
		v:   c,
		vip: v,
	}

	authInfo, err := clientImpl.login()
	if err != nil {
		return nil, fmt.Errorf("unable to login to Vault: %w", err)
	}

	clientImpl.authInfo = authInfo

	go clientImpl.renewAuthInfo()

	return clientImpl, nil
}

func (c *appRoleClient) Client() *vault.Client {
	return c.v
}

func (c *appRoleClient) login() (*vault.Secret, error) {
	vip := c.vip
	approleSecretID := &approle.SecretID{
		FromString: vip.GetString("vault.app_role_secret_id"),
	}

	// Authenticate with Vault with the AppRole auth method
	appRoleAuth, err := approle.NewAppRoleAuth(
		vip.GetString("vault.app_role_id"),
		approleSecretID,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create AppRole auth: %w", err)
	}

	authInfo, err := c.v.Auth().Login(context.Background(), appRoleAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate with Vault: %w", err)
	}
	if authInfo == nil {
		return nil, errors.New("authentication with Vault failed")
	}

	return authInfo, nil
}

func (c *appRoleClient) renewAuthInfo() {
	err := RenewLease(context.Background(), c, "auth", c.authInfo, func() (*vault.Secret, error) {
		authInfo, err := c.login()
		if err != nil {
			return nil, fmt.Errorf("unable to renew auth info: %w", err)
		}

		c.authInfo = authInfo

		return authInfo, nil
	})
	if err != nil {
		slog.Error("unable to renew auth info", slog.String(logging.KeyError, err.Error()))
		os.Exit(1)
	}
}

// SetKvSecretV2 sets the secrets in the Vault server.
func (c *appRoleClient) SetKvSecretV2(ctx context.Context, mount, name string, data map[string]any) error {
	_, err := c.v.KVv2(mount).Put(ctx, name, data)
	if err != nil {
		return fmt.Errorf("unable to set secret: %w", err)
	}

	return nil
}

// GetKvSecretV2 gets the secrets from the Vault server.
func (c *appRoleClient) GetKvSecretV2(ctx context.Context, mount, name string) (*vault.KVSecret, error) {
	secret, err := c.v.KVv2(mount).Get(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("unable to read secret: %w", err)
	} else if secret == nil {
		return nil, ErrSecretNotFound
	}

	return secret, nil
}

// GetSecret gets the secrets from the Vault server.
func (c *appRoleClient) GetSecret(ctx context.Context, path string) (*Secrets, error) {
	secret, err := c.v.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("unable to read secrets: %w", err)
	} else if secret == nil {
		return nil, ErrSecretNotFound
	}
	return &Secrets{secret}, nil
}

func (c *appRoleClient) TransitEncrypt(ctx context.Context, data string) (*Secrets, error) {
	return transitEncrypt(ctx, c, c.vip.GetString("vault.transit.path_encrypt"), data)
}

func (c *appRoleClient) TransitDecrypt(ctx context.Context, data string) (string, error) {
	return transitDecrypt(ctx, c, c.vip.GetString("vault.transit.path_decrypt"), data)
}
