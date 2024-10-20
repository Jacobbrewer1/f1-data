package vaulty

import hashiVault "github.com/hashicorp/vault/api"

func CipherTextFromSecret(transitEncryptSecret *hashiVault.Secret) string {
	if transitEncryptSecret == nil {
		return ""
	} else if transitEncryptSecret.Data == nil {
		return ""
	} else if transitEncryptSecret.Data[TransitKeyCipherText] == nil {
		return ""
	}

	ct, ok := transitEncryptSecret.Data[TransitKeyCipherText].(string)
	if !ok {
		return ""
	}

	return ct
}
