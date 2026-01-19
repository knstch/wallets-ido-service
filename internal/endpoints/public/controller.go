package public

import (
	"net/http"

	"github.com/knstch/knstch-libs/middleware"
	"github.com/knstch/knstch-libs/transport"

	"wallets-service/config"
	"wallets-service/internal/wallets"

	"github.com/knstch/knstch-libs/log"

	httptransport "github.com/go-kit/kit/transport/http"
	public "github.com/knstch/wallets-ido-api/public"

	"github.com/knstch/knstch-libs/endpoints"
)

type Controller struct {
	svc wallets.Service
	lg  *log.Logger
	cfg *config.Config

	public.UnimplementedWalletsServer
}

func NewController(svc wallets.Service, lg *log.Logger, cfg *config.Config) *Controller {
	return &Controller{
		svc: svc,
		cfg: cfg,
		lg:  lg,
	}
}

func (c *Controller) Endpoints() []endpoints.Endpoint {
	defaultMiddlewares := []middleware.Middleware{middleware.WithCookieAuth(c.cfg.JwtSecret)}

	return []endpoints.Endpoint{
		{
			Method:  http.MethodPost,
			Path:    "/addWallet",
			Handler: MakeAddWalletEndpoint(c),
			Decoder: transport.DecodeJSONRequest[public.AddWalletRequest],
			Encoder: httptransport.EncodeJSONResponse,
			Mdw:     defaultMiddlewares,
		},
		{
			Method:  http.MethodPost,
			Path:    "/verifyWallet",
			Handler: MakeVerifyWalletEndpoint(c),
			Decoder: transport.DecodeJSONRequest[public.VerifyWalletRequest],
			Encoder: httptransport.EncodeJSONResponse,
			Mdw:     defaultMiddlewares,
		},
		{
			Method:  http.MethodPost,
			Path:    "/unlinkWallet",
			Handler: MakeUnlinkWalletEndpoint(c),
			Decoder: transport.DecodeJSONRequest[public.UnlinkWalletRequest],
			Encoder: httptransport.EncodeJSONResponse,
			Mdw:     defaultMiddlewares,
		},
		{
			Method:  http.MethodGet,
			Path:    "/getWallet",
			Handler: MakeGetWalletEndpoint(c),
			Decoder: transport.DecodeDefaultRequest,
			Encoder: httptransport.EncodeJSONResponse,
			Mdw:     defaultMiddlewares,
		},
	}
}
