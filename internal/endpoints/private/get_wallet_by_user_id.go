package private

import (
	"context"
	"fmt"

	"github.com/knstch/knstch-libs/tracing"
	private "github.com/knstch/wallets-ido-api/private"

	"wallets-service/internal/domain/enum"
)

func (c *Controller) GetWalletByUserID(ctx context.Context, req *private.GetWalletByUserIDRequest) (*private.GetWalletByUserIDResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "private: GetWalletByUserID")
	defer span.End()

	wallet, err := c.svc.GetWallet(ctx, uint(req.GetUserId()))
	if err != nil {
		return nil, fmt.Errorf("svc.GetWallet: %w", err)
	}

	transportProvider, err := convertSvcProviderToTransport(wallet.Provider)
	if err != nil {
		return nil, err
	}

	return &private.GetWalletByUserIDResponse{
		Id:         uint64(wallet.ID),
		Pubkey:     wallet.Pubkey,
		Provider:   transportProvider,
		IsVerified: wallet.VerifiedAt != nil,
	}, nil
}

func convertSvcProviderToTransport(provider enum.Provider) (private.Provider, error) {
	switch provider {
	case enum.ProviderPhantom:
		return private.Provider_PROVIDER_PHANTOM, nil
	default:
		return private.Provider_PROVIDER_UNDEFINED, fmt.Errorf("unknown provider: %s", provider)
	}
}
