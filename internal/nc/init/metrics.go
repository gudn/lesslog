package nc_init

import (
	"time"

	"github.com/gudn/iinit"
	"github.com/gudn/lesslog/internal/metrics"
	"github.com/gudn/lesslog/internal/nc"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	in_msgs = prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: "nats",
		Name:      "in_msg",
		Help:      "number of incomed messages",
	})
	out_msgs = prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: "nats",
		Name:      "out_msg",
		Help:      "number of outcomed messages",
	})
	reconnects = prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: "nats",
		Name:      "reconnects_total",
		Help:      "number of reconnects",
	})
	subs_count = prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: "nats",
		Name:      "subscriptions_count",
		Help:      "active number of subscriptions",
	})
)

func init() {
	p := iinit.ParallelS(
		metrics.InitMetrics,
		InitNats,
	)

	iinit.Sequential(
		p,
		iinit.Static(func() {
			if !metrics.IsEnabled() || nc.Conn == nil {
				return
			}
			prometheus.MustRegister(in_msgs, out_msgs, reconnects, subs_count)
			go func() {
				ticker := time.NewTicker(1 * time.Second)
				defer ticker.Stop()
				for {
					select {
					case <-ticker.C:
						subs_count.Set(float64(nc.Conn.NumSubscriptions()))
						stat := nc.Conn.Stats()
						in_msgs.Set(float64(stat.InMsgs))
						out_msgs.Set(float64(stat.OutMsgs))
						reconnects.Set(float64(stat.Reconnects))
					}
				}
			}()
		}),
	)
}
