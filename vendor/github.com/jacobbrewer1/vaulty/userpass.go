package vaulty

import (
	"context"
	"fmt"
	"time"

	hashiVault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/userpass"
)

func userPassLogin(client *hashiVault.Client, username, password string) (*hashiVault.Secret, error) {
	// WARNING: A plaintext password like this is obviously insecure.
	// See the hashicorp/vault-examples repo for full examples of how to securely
	// log in to Vault using various auth methods. This function is just
	// demonstrating the basic idea that a *vault.Secret is returned by
	// the login call.
	userPassAuth, err := auth.NewUserpassAuth(
		username,
		&auth.Password{
			FromString: password,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize userpass auth method: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authInfo, err := client.Auth().Login(ctx, userPassAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to login to userpass auth method: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}

	return authInfo, nil
}
