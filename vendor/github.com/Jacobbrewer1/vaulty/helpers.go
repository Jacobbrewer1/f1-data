package vaulty

import hashiVault "github.com/hashicorp/vault/api"

func GetTransitCipherText(transitEncryptSecret *hashiVault.Secret) string {
	if transitEncryptSecret == nil {
		return ""
	}

	return transitEncryptSecret.Data[TransitCipherText].(string)
}
