package vaulty

import hashiVault "github.com/hashicorp/vault/api"

func tokenLogin(v *hashiVault.Client, token string) (*hashiVault.Secret, error) {
	v.SetToken(token)
	return v.Auth().Token().LookupSelf()
}
