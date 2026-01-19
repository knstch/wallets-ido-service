package wallets

type Challenge struct {
	UserID    uint   `json:"user_id"`
	PubKey    string `json:"pubkey"`
	Provider  string `json:"provider"`
	Nonce     string `json:"nonce"`
	ExpiresAt int64  `json:"expires_at"`
}
