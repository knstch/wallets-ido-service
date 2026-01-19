package wallets_test

import (
	"context"

	"github.com/knstch/knstch-libs/svcerrs"

	"wallets-service/internal/domain/enum"
)

func (s *WalletsServiceTestSuite) TestUnlinkWallet_HappyPath() {
	t := s.Require()
	pubkey, _ := mustGenerateSolanaKeypair(t)
	_, err := s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)

	w, err := s.svc.GetWallet(context.Background(), 1)
	t.NoError(err)

	t.NoError(s.svc.UnlinkWallet(context.Background(), w.ID, 1))

	_, err = s.svc.GetWallet(context.Background(), 1)
	requireSvcErrIs(s.T(), err, svcerrs.ErrDataNotFound)
}

func (s *WalletsServiceTestSuite) TestUnlinkWallet_WrongUser_NotFound() {
	t := s.Require()
	pubkey, _ := mustGenerateSolanaKeypair(t)
	_, err := s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)

	w, err := s.svc.GetWallet(context.Background(), 1)
	t.NoError(err)

	err = s.svc.UnlinkWallet(context.Background(), w.ID, 2)
	requireSvcErrIs(s.T(), err, svcerrs.ErrDataNotFound)
}

func (s *WalletsServiceTestSuite) TestUnlinkWallet_IdempotentSecondDelete_NotFound() {
	t := s.Require()
	pubkey, _ := mustGenerateSolanaKeypair(t)
	_, err := s.svc.AddWallet(context.Background(), 1, pubkey, enum.ProviderPhantom)
	t.NoError(err)

	w, err := s.svc.GetWallet(context.Background(), 1)
	t.NoError(err)

	t.NoError(s.svc.UnlinkWallet(context.Background(), w.ID, 1))
	err = s.svc.UnlinkWallet(context.Background(), w.ID, 1)
	requireSvcErrIs(s.T(), err, svcerrs.ErrDataNotFound)
}

