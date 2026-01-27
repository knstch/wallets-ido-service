package private

import (
	private "github.com/knstch/wallets-ido-api/private"

	"github.com/knstch/knstch-libs/log"

	"wallets-service/config"
	"wallets-service/internal/wallets"
)

type Controller struct {
	svc wallets.Service
	lg  *log.Logger
	cfg *config.Config

	private.UnimplementedWalletsPrivateServer
}

func NewController(svc wallets.Service, lg *log.Logger, cfg *config.Config) *Controller {
	return &Controller{
		svc: svc,
		cfg: cfg,
		lg:  lg,
	}
}
