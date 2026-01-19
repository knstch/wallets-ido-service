package wallets_test

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/knstch/knstch-libs/svcerrs"

	"wallets-service/internal/domain/enum"
	"wallets-service/internal/wallets"
)

func (s *WalletsServiceTestSuite) TestAddWallet_HappyPath() {
	t := s.Require()
	pubkey, _ := mustGenerateSolanaKeypair(t)
	res, err := s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)
	t.NotEmpty(res.ChallengeID)
	t.NotEmpty(res.MessageToSign)

	// Challenge should be present in Redis.
	val, err := s.rdb.Get(context.Background(), wallets.GetChallengeByIDKey(res.ChallengeID)).Result()
	t.NoError(err)
	t.NotEmpty(val)
}

func (s *WalletsServiceTestSuite) TestAddWallet_ConflictUnverified_AllowsReverify() {
	t := s.Require()
	pubkey, _ := mustGenerateSolanaKeypair(t)
	_, err := s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)

	// Second call should succeed (wallet exists but is not verified yet).
	_, err = s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)
}

func (s *WalletsServiceTestSuite) TestAddWallet_ConflictOtherUser_ReturnsConflict() {
	t := s.Require()
	pubkey, _ := mustGenerateSolanaKeypair(t)

	_, err := s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)

	_, err = s.svc.AddWallet(context.Background(), 2, pubkey, enum.ProviderPhantom)
	requireSvcErrIs(s.T(), err, svcerrs.ErrConflict)
}

func (s *WalletsServiceTestSuite) TestAddWallet_RedisDown_ReturnsErrorAndDoesNotCreateWallet() {
	t := s.Require()
	pubkey, _ := mustGenerateSolanaKeypair(t)

	// Broken redis client (dial will fail on operations).
	badRedis := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
	svc := wallets.NewService(s.logger, s.dbRepo, s.cfg, badRedis)

	_, err := svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.Error(err)

	// Ensure wallet wasn't created (DB transaction should roll back on redis.Set error).
	_, err = s.svc.GetWallet(context.Background(), 1)
	requireSvcErrIs(s.T(), err, svcerrs.ErrDataNotFound)
}

func (s *WalletsServiceTestSuite) TestAddWallet_ConflictVerified_ReturnsConflict() {
	t := s.Require()
	pubkey, priv := mustGenerateSolanaKeypair(t)

	ch, err := s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)
	sig := mustSignBase64(priv, ch.MessageToSign)
	t.NoError(s.svc.VerifyWallet(context.Background(), 1, ch.ChallengeID, sig, pubkey))

	_, err = s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	requireSvcErrIs(s.T(), err, svcerrs.ErrConflict)
}
