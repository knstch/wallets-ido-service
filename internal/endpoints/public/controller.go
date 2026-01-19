package public

import (
	"github.com/knstch/knstch-libs/middleware"

	"wallets-service/config"
	"wallets-service/internal/wallets"

	"github.com/knstch/knstch-libs/log"

	public "github.com/knstch/users-ido-api/public"

	"github.com/knstch/knstch-libs/endpoints"
)

type Controller struct {
	svc wallets.Service
	lg  *log.Logger
	cfg *config.Config

	public.UnimplementedUsersServer
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

	return []endpoints.Endpoint{}
}
