package wallets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/knstch/knstch-libs/svcerrs"
	"github.com/knstch/knstch-libs/tracing"

	"wallets-service/internal/domain/dto"
	"wallets-service/internal/domain/enum"
	"wallets-service/internal/wallets/filters"
	"wallets-service/internal/wallets/repo"
	"wallets-service/internal/wallets/utils"
)

const (
	nonceLen = 32

	challengeExpirationPeriod = time.Minute * 15

	messageToSignToVerifyWallet = "Please, verify your wallet"
)

func buildMessageToSign(challengeID, pubkey, nonce string, expiresAt int64) string {
	return fmt.Sprintf(
		"%s\n\nPubkey: %s\nChallengeId: %s\nNonce: %s\nExpiresAt: %d",
		messageToSignToVerifyWallet,
		pubkey,
		challengeID,
		nonce,
		expiresAt,
	)
}

// AddWallet creates a wallet record for the user and returns a verification challenge.
//
// Behavior:
//   - If the wallet is new, it is created in Postgres and the challenge is stored in Redis.
//   - If the wallet already exists and is NOT verified, the service re-issues a new challenge (re-verify flow).
//   - If the wallet already exists and is verified (or belongs to another user due to unique constraints),
//     AddWallet returns svcerrs.ErrConflict.
//
// The returned MessageToSign must be signed by the wallet owner and then validated via VerifyWallet.
func (s *ServiceImpl) AddWallet(ctx context.Context, userID uint, pubkey string, provider enum.Provider) (dto.ChallengeForUser, error) {
	ctx, span := tracing.StartSpan(ctx, "wallets: AddWallet")
	defer span.End()

	challengeID := uuid.New()
	expiresAt := time.Now().Add(challengeExpirationPeriod)
	nonce, err := utils.RandomString(nonceLen)
	if err != nil {
		return dto.ChallengeForUser{}, fmt.Errorf("utils.RandomString: %w", err)
	}

	challenge := &Challenge{
		UserID:    userID,
		PubKey:    pubkey,
		Provider:  provider.String(),
		Nonce:     nonce,
		ExpiresAt: expiresAt.Unix(),
	}
	msg := buildMessageToSign(challengeID.String(), pubkey, nonce, expiresAt.Unix())

	jsonChallenge, err := json.Marshal(challenge)
	if err != nil {
		return dto.ChallengeForUser{}, fmt.Errorf("json.Marshal: %w", err)
	}

	if err = s.repo.Transaction(func(st repo.Repository) error {
		if err := st.CreateWallet(ctx, userID, pubkey, provider); err != nil {
			// Bubble up conflict as-is so we can handle it outside the transaction.
			if errors.Is(err, svcerrs.ErrConflict) {
				return svcerrs.ErrConflict
			}
			return fmt.Errorf("st.CreateWallet: %w", err)
		}

		if err := s.redis.Set(
			ctx,
			GetChallengeByIDKey(challengeID.String()),
			jsonChallenge,
			challengeExpirationPeriod,
		).Err(); err != nil {
			return fmt.Errorf("redis.Set: %w", err)
		}

		return nil
	}); err != nil {
		if errors.Is(err, svcerrs.ErrConflict) {
			isVerified := false
			if _, getErr := s.repo.GetWallet(ctx, filters.WalletsFilter{
				UserID:     userID,
				Pubkey:     pubkey,
				Provider:   provider,
				IsVerified: &isVerified,
			}); getErr != nil {
				if errors.Is(getErr, svcerrs.ErrDataNotFound) {
					return dto.ChallengeForUser{}, fmt.Errorf("attempt to re-verify wallet: %w", svcerrs.ErrConflict)
				}
				return dto.ChallengeForUser{}, fmt.Errorf("repo.GetWallet: %w", getErr)
			}

			if err := s.redis.Set(
				ctx,
				GetChallengeByIDKey(challengeID.String()),
				jsonChallenge,
				challengeExpirationPeriod,
			).Err(); err != nil {
				return dto.ChallengeForUser{}, fmt.Errorf("redis.Set: %w", err)
			}
		} else {
			return dto.ChallengeForUser{}, fmt.Errorf("repo.Transaction: %w", err)
		}
	}

	return dto.ChallengeForUser{
		ChallengeID:   challengeID.String(),
		MessageToSign: msg,
	}, nil
}
