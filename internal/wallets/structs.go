package wallets

type Challenge struct {
	// UserID is the owner of the wallet being verified.
	UserID    uint   `json:"user_id"`
	// PubKey is the wallet public key (base58 string for Solana).
	PubKey    string `json:"pubkey"`
	// Provider identifies the wallet provider (e.g. "phantom").
	Provider  string `json:"provider"`
	// Nonce is a random string to prevent replay/signature reuse.
	Nonce     string `json:"nonce"`
	// ExpiresAt is a unix timestamp (seconds) after which the challenge is invalid.
	ExpiresAt int64  `json:"expires_at"`
}
