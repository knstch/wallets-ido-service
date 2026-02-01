package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	registerWalletsOnce sync.Once

	walletsAddTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "wallets_service",
		Subsystem: "wallets",
		Name:      "add_total",
		Help:      "Total number of AddWallets calls.",
	})
)

func registerWallets() {
	registerWalletsOnce.Do(func() {
		prometheus.MustRegister(
			walletsAddTotal,
		)
	})
}

// IncAddWallet increments the AddWallets Prometheus counter.
func IncAddWallet() {
	registerWallets()
	walletsAddTotal.Inc()
}
