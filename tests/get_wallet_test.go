package wallets_test

import (
	"context"

	"github.com/knstch/knstch-libs/svcerrs"
)

func (s *WalletsServiceTestSuite) TestGetWallet_NotFound() {
	_, err := s.svc.GetWallet(context.Background(), 1)
	requireSvcErrIs(s.T(), err, svcerrs.ErrDataNotFound)
}

