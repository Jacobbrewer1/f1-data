package vaulty

import (
	"context"
	"fmt"
	"time"

	hashiVault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"
)

func kubernetesLogin(client *hashiVault.Client, role string, token auth.LoginOption) (*hashiVault.Secret, error) {
	k8sAuth, err := auth.NewKubernetesAuth(
		role,
		token,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize kubernetes auth method: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authInfo, err := client.Auth().Login(ctx, k8sAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to login to kubernetes auth method: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}

	client.SetToken(authInfo.Auth.ClientToken)

	return authInfo, nil
}
