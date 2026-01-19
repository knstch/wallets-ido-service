package wallets_test

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/knstch/knstch-libs/svcerrs"

	"wallets-service/internal/domain/enum"
	"wallets-service/internal/wallets"
)

func (s *WalletsServiceTestSuite) TestVerifyWallet_HappyPath() {
	t := s.Require()
	pubkey, priv := mustGenerateSolanaKeypair(t)

	ch, err := s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)

	sig := mustSignBase64(priv, ch.MessageToSign)
	t.NoError(s.svc.VerifyWallet(context.Background(), 1, ch.ChallengeID, sig, pubkey))

	// After success, challenge should be gone (best effort).
	_, err = s.rdb.Get(context.Background(), wallets.GetChallengeByIDKey(ch.ChallengeID)).Result()
	t.ErrorIs(err, redis.Nil)

	// Wallet should be verified in DB.
	w, err := s.svc.GetWallet(context.Background(), 1)
	t.NoError(err)
	t.NotNil(w.VerifiedAt)
}

func (s *WalletsServiceTestSuite) TestVerifyWallet_ChallengeNotFound() {
	err := s.svc.VerifyWallet(context.Background(), 1, "missing", "sig", "pub")
	requireSvcErrIs(s.T(), err, svcerrs.ErrDataNotFound)
}

func (s *WalletsServiceTestSuite) TestVerifyWallet_RedisDown_ReturnsError() {
	t := s.Require()
	pubkey, _ := mustGenerateSolanaKeypair(t)
	ch, err := s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)

	badRedis := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
	svc := wallets.NewService(s.logger, s.dbRepo, s.cfg, badRedis)

	err = svc.VerifyWallet(context.Background(), 1, ch.ChallengeID, "sig", pubkey)
	t.Error(err)
}

func (s *WalletsServiceTestSuite) TestVerifyWallet_UserMismatch_NotFound() {
	t := s.Require()
	pubkey, _ := mustGenerateSolanaKeypair(t)
	ch, err := s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)

	err = s.svc.VerifyWallet(context.Background(), 2, ch.ChallengeID, "sig", pubkey)
	requireSvcErrIs(s.T(), err, svcerrs.ErrDataNotFound)
}

func (s *WalletsServiceTestSuite) TestVerifyWallet_PubkeyMismatch_InvalidData() {
	t := s.Require()
	pubkey, _ := mustGenerateSolanaKeypair(t)
	ch, err := s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)

	err = s.svc.VerifyWallet(context.Background(), 1, ch.ChallengeID, "sig", "another-pubkey")
	requireSvcErrIs(s.T(), err, svcerrs.ErrInvalidData)
}

func (s *WalletsServiceTestSuite) TestVerifyWallet_MalformedChallengeJSON_ReturnsError() {
	t := s.Require()
	// Put garbage into Redis.
	t.NoError(s.rdb.Set(context.Background(), wallets.GetChallengeByIDKey("bad-json"), "{", time.Minute).Err())

	err := s.svc.VerifyWallet(context.Background(), 1, "bad-json", "sig", "pub")
	t.Error(err)
}

func (s *WalletsServiceTestSuite) TestVerifyWallet_InvalidChallengePubkey_InvalidData() {
	t := s.Require()
	// Store a challenge with invalid base58 pubkey; signature doesn't matter (will fail before verify).
	bad := wallets.Challenge{
		UserID:    1,
		PubKey:    "not-base58!!!",
		Provider:  enum.ProviderPhantom.String(),
		Nonce:     "n",
		ExpiresAt: time.Now().Add(time.Minute).Unix(),
	}
	raw, err := json.Marshal(&bad)
	t.NoError(err)
	t.NoError(s.rdb.Set(context.Background(), wallets.GetChallengeByIDKey("bad-pk"), raw, time.Minute).Err())

	err = s.svc.VerifyWallet(context.Background(), 1, "bad-pk", "sig", bad.PubKey)
	requireSvcErrIs(s.T(), err, svcerrs.ErrInvalidData)
}

func (s *WalletsServiceTestSuite) TestVerifyWallet_InvalidSignatureEncoding_InvalidData() {
	t := s.Require()
	pubkey, _ := mustGenerateSolanaKeypair(t)
	ch, err := s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)

	// Not base64.
	err = s.svc.VerifyWallet(context.Background(), 1, ch.ChallengeID, "%%%notbase64%%%", pubkey)
	requireSvcErrIs(s.T(), err, svcerrs.ErrInvalidData)
}

func (s *WalletsServiceTestSuite) TestVerifyWallet_InvalidSignature_InvalidData() {
	t := s.Require()
	pubkey, priv := mustGenerateSolanaKeypair(t)
	ch, err := s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)

	sig := mustSignBase64(priv, ch.MessageToSign+"tampered")
	err = s.svc.VerifyWallet(context.Background(), 1, ch.ChallengeID, sig, pubkey)
	requireSvcErrIs(s.T(), err, svcerrs.ErrInvalidData)
}

func (s *WalletsServiceTestSuite) TestVerifyWallet_ExpiredChallenge_NotFound() {
	t := s.Require()
	pubkey, priv := mustGenerateSolanaKeypair(t)
	ch, err := s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)

	// Force expires_at into the past in Redis.
	raw, err := s.rdb.Get(context.Background(), wallets.GetChallengeByIDKey(ch.ChallengeID)).Result()
	t.NoError(err)

	var stored wallets.Challenge
	t.NoError(json.Unmarshal([]byte(raw), &stored))
	stored.ExpiresAt = time.Now().Add(-time.Second).Unix()
	b, err := json.Marshal(&stored)
	t.NoError(err)
	t.NoError(s.rdb.Set(context.Background(), wallets.GetChallengeByIDKey(ch.ChallengeID), b, time.Minute).Err())

	sig := mustSignBase64(priv, ch.MessageToSign)
	err = s.svc.VerifyWallet(context.Background(), 1, ch.ChallengeID, sig, pubkey)
	requireSvcErrIs(s.T(), err, svcerrs.ErrDataNotFound)
}

func (s *WalletsServiceTestSuite) TestVerifyWallet_InvalidProviderInChallenge_ReturnsError() {
	t := s.Require()
	pubkey, priv := mustGenerateSolanaKeypair(t)
	ch, err := s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)

	// Mutate provider stored in Redis to an unknown one.
	raw, err := s.rdb.Get(context.Background(), wallets.GetChallengeByIDKey(ch.ChallengeID)).Result()
	t.NoError(err)
	var stored wallets.Challenge
	t.NoError(json.Unmarshal([]byte(raw), &stored))
	stored.Provider = "unknown-provider"
	b, err := json.Marshal(&stored)
	t.NoError(err)
	t.NoError(s.rdb.Set(context.Background(), wallets.GetChallengeByIDKey(ch.ChallengeID), b, time.Minute).Err())

	sig := mustSignBase64(priv, ch.MessageToSign)
	err = s.svc.VerifyWallet(context.Background(), 1, ch.ChallengeID, sig, pubkey)
	t.Error(err)
}

func (s *WalletsServiceTestSuite) TestVerifyWallet_WalletDeletedBeforeVerify_NotFound() {
	t := s.Require()
	pubkey, priv := mustGenerateSolanaKeypair(t)
	ch, err := s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)

	w, err := s.svc.GetWallet(context.Background(), 1)
	t.NoError(err)
	t.NoError(s.svc.UnlinkWallet(context.Background(), w.ID, 1))

	sig := mustSignBase64(priv, ch.MessageToSign)
	err = s.svc.VerifyWallet(context.Background(), 1, ch.ChallengeID, sig, pubkey)
	requireSvcErrIs(s.T(), err, svcerrs.ErrDataNotFound)
}

func (s *WalletsServiceTestSuite) TestVerifyWallet_Replay_ReturnsNotFound() {
	t := s.Require()
	pubkey, priv := mustGenerateSolanaKeypair(t)
	ch, err := s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)

	sig := mustSignBase64(priv, ch.MessageToSign)
	t.NoError(s.svc.VerifyWallet(context.Background(), 1, ch.ChallengeID, sig, pubkey))

	// Second try with same challengeID should fail either at Redis (deleted) or at DB (already verified).
	err = s.svc.VerifyWallet(context.Background(), 1, ch.ChallengeID, sig, pubkey)
	requireSvcErrIs(s.T(), err, svcerrs.ErrDataNotFound)
}
