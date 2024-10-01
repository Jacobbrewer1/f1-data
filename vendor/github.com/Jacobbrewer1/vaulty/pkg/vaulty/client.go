package vaulty

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/Jacobbrewer1/vaulty/pkg/logging"
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

	// GetKvSecretV2 returns a map of secrets for the given path.
	GetKvSecretV2(ctx context.Context, name string) (*hashiVault.KVSecret, error)

	// GetSecret returns a map of secrets for the given path.
	GetSecret(ctx context.Context, path string) (*hashiVault.Secret, error)

	// TransitEncrypt encrypts the given data.
	TransitEncrypt(ctx context.Context, data string) (*hashiVault.Secret, error)

	// TransitDecrypt decrypts the given data.
	TransitDecrypt(ctx context.Context, data string) (string, error)
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
		slog.Error("unable to renew auth info", slog.String(logging.KeyError, err.Error()))
		os.Exit(1)
	}
}

func (c *client) Client() *hashiVault.Client {
	return c.v
}

func (c *client) GetKvSecretV2(ctx context.Context, name string) (*hashiVault.KVSecret, error) {
	secret, err := c.v.KVv2(c.kvv2Mount).Get(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("unable to read secret: %w", err)
	} else if secret == nil {
		return nil, ErrSecretNotFound
	}
	return secret, nil
}

func (c *client) GetSecret(ctx context.Context, path string) (*hashiVault.Secret, error) {
	secret, err := c.v.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("unable to read secrets: %w", err)
	} else if secret == nil {
		return nil, ErrSecretNotFound
	}
	return secret, nil
}

func (c *client) TransitEncrypt(ctx context.Context, data string) (*hashiVault.Secret, error) {
	plaintext := base64.StdEncoding.EncodeToString([]byte(data))

	// Encrypt the data using the transit engine
	encryptData, err := c.v.Logical().WriteWithContext(ctx, c.transitPathEncrypt, map[string]any{
		"plaintext": plaintext,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt data: %w", err)
	}

	return encryptData, nil
}

func (c *client) TransitDecrypt(ctx context.Context, data string) (string, error) {
	// Decrypt the data using the transit engine
	decryptData, err := c.v.Logical().WriteWithContext(ctx, c.transitPathDecrypt, map[string]any{
		"ciphertext": data,
	})
	if err != nil {
		return "", fmt.Errorf("unable to decrypt data: %w", err)
	}

	// Decode the base64 encoded data
	decodedData, err := base64.StdEncoding.DecodeString(decryptData.Data["plaintext"].(string))
	if err != nil {
		return "", fmt.Errorf("unable to decode data: %w", err)
	}

	return string(decodedData), nil
}
