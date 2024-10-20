package vaulty

import (
	"context"
	"encoding/base64"
	"fmt"

	hashiVault "github.com/hashicorp/vault/api"
)

type SecretPath struct {
	r       Client
	mount   string
	prefix  string
	name    string
	version int
}

func (c *SecretPath) path() string {
	if c.prefix != "" {
		return fmt.Sprintf("%s/%s", c.prefix, c.name)
	}

	return c.name
}

func (c *SecretPath) pathWithType(k string) string {
	if c.prefix != "" {
		return fmt.Sprintf("%s/%s/%s", c.prefix, k, c.name)
	}

	return fmt.Sprintf("%s/%s", k, c.name)
}

func (c *SecretPath) GetKvSecretV2(ctx context.Context) (*hashiVault.KVSecret, error) {
	secret, err := c.r.Client().KVv2(c.mount).GetVersion(ctx, c.path(), c.version)
	if err != nil {
		return nil, fmt.Errorf("unable to read secret: %w", err)
	} else if secret == nil {
		return nil, ErrSecretNotFound
	}
	return secret, nil
}

func (c *SecretPath) GetSecret(ctx context.Context) (*hashiVault.Secret, error) {
	secret, err := c.r.Client().Logical().ReadWithContext(ctx, c.path())
	if err != nil {
		return nil, fmt.Errorf("unable to read secrets: %w", err)
	} else if secret == nil {
		return nil, ErrSecretNotFound
	}
	return secret, nil
}

func (c *SecretPath) TransitEncrypt(ctx context.Context, data string) (*hashiVault.Secret, error) {
	plaintext := base64.StdEncoding.EncodeToString([]byte(data))

	// Encrypt the data using the transit engine
	encryptData, err := c.r.Client().Logical().WriteWithContext(ctx, c.pathWithType(pathKeyTransitEncrypt), map[string]any{
		TransitKeyPlainText: plaintext,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt data: %w", err)
	}

	return encryptData, nil
}

func (c *SecretPath) TransitDecrypt(ctx context.Context, data string) (string, error) {
	// Decrypt the data using the transit engine
	decryptData, err := c.r.Client().Logical().WriteWithContext(ctx, c.pathWithType(pathKeyTransitDecrypt), map[string]any{
		TransitKeyCipherText: data,
	})
	if err != nil {
		return "", fmt.Errorf("unable to decrypt data: %w", err)
	}

	// Decode the base64 encoded data
	decodedData, err := base64.StdEncoding.DecodeString(decryptData.Data[TransitKeyPlainText].(string))
	if err != nil {
		return "", fmt.Errorf("unable to decode data: %w", err)
	}

	return string(decodedData), nil
}
