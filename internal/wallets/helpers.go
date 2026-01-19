package wallets

// GetChallengeByIDKey builds the Redis key used to store a challenge by its ID.
func GetChallengeByIDKey(id string) string {
	return "challenge:" + id
}
