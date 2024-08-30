package vault

import (
	"context"
	"encoding/base64"
	"fmt"
)

// transitEncrypt encrypts the given data.
func transitEncrypt(ctx context.Context, client Client, path, data string) (*Secrets, error) {
	plaintext := base64.StdEncoding.EncodeToString([]byte(data))

	// Encrypt the data using the transit engine
	encryptData, err := client.Client().Logical().WriteWithContext(ctx, path, map[string]any{
		"plaintext": plaintext,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt data: %w", err)
	}

	return &Secrets{encryptData}, nil
}

// transitDecrypt decrypts the given data.
func transitDecrypt(ctx context.Context, client Client, path, data string) (string, error) {
	// Decrypt the data using the transit engine
	decryptData, err := client.Client().Logical().WriteWithContext(ctx, path, map[string]any{
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
