package db_init

import (
	"time"

	"github.com/gudn/iinit"
	"github.com/gudn/lesslog/internal/db"
	"github.com/gudn/lesslog/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	acquited_count = prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: "pg",
		Name: "acquited_count",
		Help: "number of currently acquired connections in the pool",
	})
	idle_count = prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: "pg",
		Name: "idle_count",
		Help: "number of currently idle conns in the pool",
	})
	total_count = prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: "pg",
		Name: "total_count",
		Help: "total number of resources currently in the pool",
	})
)

func init() {
	p := iinit.ParallelS(
		metrics.InitMetrics,
		InitDb,
	)

	iinit.Sequential(
		p,
		iinit.Static(func() {
			if !metrics.IsEnabled() || db.Pool == nil {
				return
			}
			prometheus.MustRegister(acquited_count, idle_count, total_count)
			go func() {
				ticker := time.NewTicker(1 * time.Second)
				defer ticker.Stop()
				for {
					select {
					case <-ticker.C:
						stat := db.Pool.Stat()
						acquited_count.Set(float64(stat.AcquiredConns()))
						idle_count.Set(float64(stat.IdleConns()))
						total_count.Set(float64(stat.TotalConns()))
					case <-db.Ctx.Done():
						return
					}
				}
			}()
		}),
	)
}
