package vaulty

import (
	"context"
	"errors"
	"fmt"
	"time"

	hashiVault "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
)

func appRoleLogin(client *hashiVault.Client, roleID, secretID string) (*hashiVault.Secret, error) {
	appRoleSecretID := &approle.SecretID{
		FromString: secretID,
	}

	// Authenticate with Vault with the AppRole auth method
	appRoleAuth, err := approle.NewAppRoleAuth(
		roleID,
		appRoleSecretID,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create AppRole auth: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authInfo, err := client.Auth().Login(ctx, appRoleAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate with Vault: %w", err)
	}
	if authInfo == nil {
		return nil, errors.New("authentication with Vault failed")
	}

	return authInfo, nil
}
