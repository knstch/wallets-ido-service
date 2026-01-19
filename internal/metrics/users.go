package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	registerOnce sync.Once

	usersCreatedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "users_service",
		Subsystem: "wallets",
		Name:      "created_total",
		Help:      "Total number of newly created wallets",
	})
)

func register() {
	registerOnce.Do(func() {
		prometheus.MustRegister(usersCreatedTotal)
	})
}

// IncUsersCreated increments the "wallets created" Prometheus counter.
func IncUsersCreated() {
	register()
	usersCreatedTotal.Inc()
}
