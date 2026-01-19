package public

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/knstch-libs/auth"
	"github.com/knstch/knstch-libs/tracing"
	public "github.com/knstch/wallets-ido-api/public"

	"wallets-service/internal/domain/enum"
)

func MakeGetWalletEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.GetWallet(ctx, request.(*public.GetWalletRequest))
	}
}

func (c *Controller) GetWallet(ctx context.Context, _ *public.GetWalletRequest) (*public.GetWalletResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "public: GetWallet")
	defer span.End()

	user, err := auth.GetUserData(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserData: %w", err)
	}

	wallet, err := c.svc.GetWallet(ctx, user.UserID)
	if err != nil {
		return nil, fmt.Errorf("svc.GetWallet: %w", err)
	}

	transportProvider, err := convertSvcProviderToTransport(wallet.Provider)
	if err != nil {
		return nil, err
	}

	return &public.GetWalletResponse{
		Id:         uint64(wallet.ID),
		Provider:   transportProvider,
		IsVerified: wallet.VerifiedAt != nil,
	}, nil
}

func convertSvcProviderToTransport(provider enum.Provider) (public.Provider, error) {
	switch provider {
	case enum.ProviderPhantom:
		return public.Provider_PROVIDER_PHANTOM, nil
	default:
		return public.Provider_PROVIDER_UNDEFINED, fmt.Errorf("unknown provider: %s", provider)
	}
}
