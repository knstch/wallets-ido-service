package public

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/knstch-libs/auth"
	"github.com/knstch/knstch-libs/tracing"
	public "github.com/knstch/wallets-ido-api/public"
)

func MakeVerifyWalletEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.VerifyWallet(ctx, request.(*public.VerifyWalletRequest))
	}
}

func (c *Controller) VerifyWallet(ctx context.Context, req *public.VerifyWalletRequest) (*public.VerifyWalletResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "public: VerifyWallet")
	defer span.End()

	user, err := auth.GetUserData(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserData: %w", err)
	}

	if err = c.svc.VerifyWallet(ctx, user.UserID, req.GetChallengeId(), req.GetSignature(), req.GetPubkey()); err != nil {
		return nil, fmt.Errorf("svc.VerifyWallet: %w", err)
	}

	return &public.VerifyWalletResponse{}, nil
}
