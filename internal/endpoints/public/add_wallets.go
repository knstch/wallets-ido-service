package public

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/knstch-libs/auth"
	"github.com/knstch/knstch-libs/tracing"
	public "github.com/knstch/wallets-ido-api/public"
)

func MakeAddWalletsEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.AddWallets(ctx, request.(*public.AddWalletsRequest))
	}
}

func (c *Controller) AddWallets(ctx context.Context, req *public.AddWalletsRequest) (*public.AddWalletsResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "public: AddWallets")
	defer span.End()

	user, err := auth.GetUserData(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserData: %w", err)
	}

	pubkeys := req.GetPubkeys()
	if err := c.svc.AddWallets(ctx, user.UserID, pubkeys); err != nil {
		return nil, fmt.Errorf("svc.AddWallets: %w", err)
	}

	return &public.AddWalletsResponse{}, nil
}

