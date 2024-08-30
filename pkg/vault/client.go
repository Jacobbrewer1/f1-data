package vault

import (
	"context"

	vault "github.com/hashicorp/vault/api"
)

var (
	ErrSecretNotFound = vault.ErrSecretNotFound
)

type Client interface {
	Client() *vault.Client

	// SetKvSecretV2 sets a map of secrets at the given path.
	SetKvSecretV2(ctx context.Context, mount, name string, data map[string]any) error

	// GetKvSecretV2 returns a map of secrets for the given path.
	GetKvSecretV2(ctx context.Context, mount, name string) (*vault.KVSecret, error)

	// GetSecret returns a map of secrets for the given path.
	GetSecret(ctx context.Context, path string) (*Secrets, error)

	// TransitEncrypt encrypts the given data.
	TransitEncrypt(ctx context.Context, data string) (*Secrets, error)

	// TransitDecrypt decrypts the given data.
	TransitDecrypt(ctx context.Context, data string) (string, error)
}

type RenewalFunc func() (*vault.Secret, error)

type Secrets struct {
	*vault.Secret
}

func (s *Secrets) Get(key string) any {
	if s.Secret == nil {
		return nil
	} else if s.Secret.Data == nil {
		return nil
	}
	return s.Secret.Data[key]
}

// CreateMockSecret should be used for testing purposes only.
func CreateMockSecret(key string, value any) *Secrets {
	return &Secrets{
		Secret: &vault.Secret{
			Data: map[string]any{
				key: value,
			},
		},
	}
}
