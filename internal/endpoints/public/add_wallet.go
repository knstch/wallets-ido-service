package public

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/knstch-libs/auth"
	"github.com/knstch/knstch-libs/svcerrs"
	"github.com/knstch/knstch-libs/tracing"
	public "github.com/knstch/wallets-ido-api/public"

	"wallets-service/internal/domain/enum"
)

func MakeAddWalletEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.AddWallet(ctx, request.(*public.AddWalletRequest))
	}
}

func (c *Controller) AddWallet(ctx context.Context, req *public.AddWalletRequest) (*public.AddWalletResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "public: AddWallet")
	defer span.End()

	user, err := auth.GetUserData(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserData: %w", err)
	}

	provider, err := convertTransportProviderToService(req.GetProvider())
	if err != nil {
		return nil, fmt.Errorf("enum.GetProvider: %w", err)
	}

	challenge, err := c.svc.AddWallet(ctx, user.UserID, req.GetPubkey(), provider)
	if err != nil {
		return nil, fmt.Errorf("svc.AddWallet: %w", err)
	}

	return &public.AddWalletResponse{
		ChallengeId:   challenge.ChallengeID,
		MessageToSign: challenge.MessageToSign,
	}, nil
}

func convertTransportProviderToService(provider public.Provider) (enum.Provider, error) {
	switch provider {
	case public.Provider_PROVIDER_PHANTOM:
		return enum.ProviderPhantom, nil
	default:
		return "", fmt.Errorf("unknown provider %s: %w", provider, svcerrs.ErrInvalidData)
	}
}
