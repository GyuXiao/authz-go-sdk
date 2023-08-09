package sdk

type Credential struct {
	SecretID string
	SecretKey string
}

func NewCredentials(secretID, secretKey string) *Credential {
	return &Credential{
		SecretID: secretID,
		SecretKey: secretKey,
	}
}