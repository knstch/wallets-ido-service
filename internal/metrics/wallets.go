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
		Help:      "Total number of AddWallet calls.",
	})
	walletsVerifyTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "wallets_service",
		Subsystem: "wallets",
		Name:      "verify_total",
		Help:      "Total number of VerifyWallet calls.",
	})
	walletsUnlinkTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "wallets_service",
		Subsystem: "wallets",
		Name:      "unlink_total",
		Help:      "Total number of UnlinkWallet calls.",
	})
)

func registerWallets() {
	registerWalletsOnce.Do(func() {
		prometheus.MustRegister(
			walletsAddTotal,
			walletsVerifyTotal,
			walletsUnlinkTotal,
		)
	})
}

// IncAddWallet increments the AddWallet Prometheus counter.
func IncAddWallet() {
	registerWallets()
	walletsAddTotal.Inc()
}

// IncVerifyWallet increments the VerifyWallet Prometheus counter.
func IncVerifyWallet() {
	registerWallets()
	walletsVerifyTotal.Inc()
}

// IncUnlinkWallet increments the UnlinkWallet Prometheus counter.
func IncUnlinkWallet() {
	registerWallets()
	walletsUnlinkTotal.Inc()
}
