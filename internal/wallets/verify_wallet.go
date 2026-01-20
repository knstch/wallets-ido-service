package wallets

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mr-tron/base58"

	"github.com/knstch/knstch-libs/log"
	"github.com/knstch/knstch-libs/svcerrs"
	"github.com/knstch/knstch-libs/tracing"
	"github.com/redis/go-redis/v9"

	"wallets-service/internal/domain/enum"
	"wallets-service/internal/metrics"
	"wallets-service/internal/wallets/filters"
)

func verifySolanaSignMessage(pubkey string, messageToSign []byte, signature string) (bool, error) {
	if pubkey == "" {
		return false, fmt.Errorf("pubkey is empty: %w", svcerrs.ErrInvalidData)
	}
	if len(messageToSign) == 0 {
		return false, fmt.Errorf("message is empty: %w", svcerrs.ErrInvalidData)
	}
	if signature == "" {
		return false, fmt.Errorf("signature is empty: %w", svcerrs.ErrInvalidData)
	}

	pubKeyBytes, err := base58.Decode(pubkey)
	if err != nil {
		return false, fmt.Errorf("base58.Decode: %w", errors.Join(err, svcerrs.ErrInvalidData))
	}
	if len(pubKeyBytes) != ed25519.PublicKeySize {
		return false, fmt.Errorf("invalid pubkey length: got %d, want %d: %w", len(pubKeyBytes), ed25519.PublicKeySize, svcerrs.ErrInvalidData)
	}

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		sigBytes, err = base64.RawURLEncoding.DecodeString(signature)
		if err != nil {
			return false, fmt.Errorf("base64.RawURLEncoding.DecodeString: %w", errors.Join(err, svcerrs.ErrInvalidData))
		}
	}
	if len(sigBytes) != ed25519.SignatureSize {
		return false, fmt.Errorf("invalid signature length: got %d, want %d: %w", len(sigBytes), ed25519.SignatureSize, svcerrs.ErrInvalidData)
	}

	ok := ed25519.Verify(pubKeyBytes, messageToSign, sigBytes)
	return ok, nil
}

// VerifyWallet validates a user's signature for a previously issued challenge and marks the wallet as verified.
//
// VerifyWallet:
// - loads the challenge JSON from Redis by challengeID
// - validates that it belongs to the user and is not expired
// - rebuilds the original message and verifies the Solana ed25519 signature
// - marks the wallet verified in Postgres
//
// On success it attempts to delete the Redis challenge key (best-effort).
func (s *ServiceImpl) VerifyWallet(ctx context.Context, userID uint, challengeID, signature, pubkey string) error {
	defer metrics.IncVerifyWallet()

	ctx, span := tracing.StartSpan(ctx, "wallets: VerifyWallet")
	defer span.End()

	challengeFromDB, err := s.redis.Get(ctx, GetChallengeByIDKey(challengeID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("challenge not found: %w", svcerrs.ErrDataNotFound)
		}
		return fmt.Errorf("redis.Get: %w", err)
	}

	var challenge Challenge
	if err = json.Unmarshal([]byte(challengeFromDB), &challenge); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	if challenge.UserID != userID {
		return fmt.Errorf("challenge not found: %w", svcerrs.ErrDataNotFound)
	}
	if pubkey != "" && challenge.PubKey != pubkey {
		return fmt.Errorf("pubkey mismatch: %w", svcerrs.ErrInvalidData)
	}
	if challenge.ExpiresAt != 0 && time.Now().Unix() > challenge.ExpiresAt {
		return fmt.Errorf("challenge expired: %w", svcerrs.ErrDataNotFound)
	}

	msg := buildMessageToSign(challengeID, challenge.PubKey, challenge.Nonce, challenge.ExpiresAt)
	verified, err := verifySolanaSignMessage(challenge.PubKey, []byte(msg), signature)
	if err != nil {
		return fmt.Errorf("verifySolanaSignMessage: %w", err)
	}
	if !verified {
		return fmt.Errorf("solana signature is invalid: %w", svcerrs.ErrInvalidData)
	}

	provider, err := enum.GetProvider(challenge.Provider)
	if err != nil {
		return fmt.Errorf("enum.GetProvider: %w", err)
	}

	isVerifiedFilter := false
	if err = s.repo.VerifyWallet(ctx, filters.WalletsFilter{
		UserID:     userID,
		Pubkey:     challenge.PubKey,
		Provider:   provider,
		IsVerified: &isVerifiedFilter,
	}); err != nil {
		return fmt.Errorf("repo.VerifyWallet: %w", err)
	}

	// We don't want to fail verification after DB update if Redis deletion fails,
	// but we also don't want to silently ignore the error.
	if err := s.redis.Del(ctx, GetChallengeByIDKey(challengeID)).Err(); err != nil && !errors.Is(err, redis.Nil) {
		s.lg.Error("redis.Del(challenge) failed", err,
			log.AddMessage("challenge_id", challengeID),
		)
	}

	return nil
}
