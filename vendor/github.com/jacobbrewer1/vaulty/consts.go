package vaulty

const (
	loggingKeyError = "err"

	pathKeyTransitDecrypt = "decrypt"
	pathKeyTransitEncrypt = "encrypt"

	TransitKeyCipherText = "ciphertext"
	TransitKeyPlainText  = "plaintext"

	envKubernetesRole  = "KUBERNETES_ROLE"
	envKubernetesToken = "KUBERNETES_TOKEN"

	kubernetesServiceAccountTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)
