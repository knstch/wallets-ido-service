package private

import (
	"context"
	"fmt"

	"github.com/knstch/knstch-libs/tracing"
	private "github.com/knstch/wallets-ido-api/private"
)

func (c *Controller) GetWalletsByUserID(ctx context.Context, req *private.GetWalletsByUserIDRequest) (*private.GetWalletsByUserIDResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "private: GetWalletsByUserID")
	defer span.End()

	wallets, err := c.svc.GetWallets(ctx, uint(req.GetUserId()))
	if err != nil {
		return nil, fmt.Errorf("svc.GetWallets: %w", err)
	}

	result := make([]*private.Wallet, 0, len(wallets))
	for _, wallet := range wallets {
		result = append(result, &private.Wallet{
			Id:     uint64(wallet.ID),
			Pubkey: wallet.Pubkey,
		})
	}

	return &private.GetWalletsByUserIDResponse{
		Wallets: result,
	}, nil
}
