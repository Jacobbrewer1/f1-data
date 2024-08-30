package vault

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Jacobbrewer1/f1-data/pkg/logging"
	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/userpass"
	"github.com/spf13/viper"
)

type userPassClient struct {
	v        *vault.Client
	authInfo *vault.Secret
	vip      *viper.Viper
}

func NewClientUserPass(v *viper.Viper) (Client, error) {
	config := vault.DefaultConfig()
	config.Address = v.GetString("vault.address")

	c, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	clientImpl := &userPassClient{
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

func (c *userPassClient) Client() *vault.Client {
	return c.v
}

func (c *userPassClient) renewAuthInfo() {
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

func (c *userPassClient) login() (*vault.Secret, error) {
	// WARNING: A plaintext password like this is obviously insecure.
	// See the hashicorp/vault-examples repo for full examples of how to securely
	// log in to Vault using various auth methods. This function is just
	// demonstrating the basic idea that a *vault.Secret is returned by
	// the login call.
	userpassAuth, err := auth.NewUserpassAuth(c.vip.GetString("vault.auth.username"), &auth.Password{FromString: c.vip.GetString("vault.auth.password")})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize userpass auth method: %w", err)
	}

	authInfo, err := c.v.Auth().Login(context.Background(), userpassAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to login to userpass auth method: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}

	return authInfo, nil
}

// SetKvSecretV2 sets the secret at the given path in Vault.
func (c *userPassClient) SetKvSecretV2(ctx context.Context, mount, name string, data map[string]any) error {
	_, err := c.v.KVv2(mount).Put(ctx, name, data)
	if err != nil {
		return fmt.Errorf("unable to set secret: %w", err)
	}

	return nil
}

// GetKvSecretV2 retrieves the secret at the given path from Vault.
func (c *userPassClient) GetKvSecretV2(ctx context.Context, mount, name string) (*vault.KVSecret, error) {
	secret, err := c.v.KVv2(mount).Get(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("unable to read secret: %w", err)
	} else if secret == nil {
		return nil, ErrSecretNotFound
	}

	return secret, nil
}

// SetSecret sets the secret at the given path in Vault.
func (c *userPassClient) SetSecret(ctx context.Context, path string, data map[string]interface{}) error {
	_, err := c.v.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		return fmt.Errorf("unable to write secrets: %w", err)
	}
	return nil
}

// GetSecret retrieves the secret at the given path from Vault.
func (c *userPassClient) GetSecret(ctx context.Context, path string) (*Secrets, error) {
	secret, err := c.v.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("unable to read secrets: %w", err)
	} else if secret == nil {
		return nil, ErrSecretNotFound
	}
	return &Secrets{secret}, nil
}

// TransitEncrypt encrypts the given data.
func (c *userPassClient) TransitEncrypt(ctx context.Context, data string) (*Secrets, error) {
	return transitEncrypt(ctx, c, c.vip.GetString("vault.transit.path_encrypt"), data)
}

// TransitDecrypt decrypts the given data.
func (c *userPassClient) TransitDecrypt(ctx context.Context, data string) (string, error) {
	return transitDecrypt(ctx, c, c.vip.GetString("vault.transit.path_decrypt"), data)
}
