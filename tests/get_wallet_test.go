package wallets_test

import (
	"context"
)

func (s *WalletsServiceTestSuite) TestGetWallets_NotFound() {
	wallets, err := s.svc.GetWallets(context.Background(), 1)
	t := s.Require()
	t.NoError(err)
	t.Empty(wallets)
}

