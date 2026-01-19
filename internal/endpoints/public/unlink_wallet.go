package public

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/knstch-libs/auth"
	"github.com/knstch/knstch-libs/tracing"
	public "github.com/knstch/wallets-ido-api/public"
)

func MakeUnlinkWalletEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.UnlinkWallet(ctx, request.(*public.UnlinkWalletRequest))
	}
}

func (c *Controller) UnlinkWallet(ctx context.Context, req *public.UnlinkWalletRequest) (*public.UnlinkWalletResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "public: UnlinkWallet")
	defer span.End()

	user, err := auth.GetUserData(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserData: %w", err)
	}

	if err = c.svc.UnlinkWallet(ctx, uint(req.GetWalletId()), user.UserID); err != nil {
		return nil, fmt.Errorf("svc.UnlinkWallet: %w", err)
	}

	return &public.UnlinkWalletResponse{}, nil
}
