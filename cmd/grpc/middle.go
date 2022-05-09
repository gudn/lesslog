package main

import (
	"context"
	"errors"

	"github.com/gudn/lesslog/internal/logging"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var (
	requests_total = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "grpc",
			Name:      "requests_total",
			Help:      "total number of processed requests",
		},
		[]string{"method", "success"},
	)
	requests_in_flight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "grpc",
			Name:      "requests_in_flight",
			Help:      "number of inprocessing reqests",
		},
		[]string{"method"},
	)
)

func streamMiddle(
	srv any,
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	requests_in_flight.
		With(prometheus.Labels{"method": info.FullMethod}).
		Inc()
	err := handler(srv, ss)
	if errors.Is(err, context.Canceled) {
		err = nil
	}
	logging.LogRequest(info.FullMethod, err)
	succ := "yes"
	if err != nil {
		succ = "no"
	}
	requests_in_flight.
		With(prometheus.Labels{"method": info.FullMethod}).
		Dec()
	requests_total.With(
		prometheus.Labels{
			"method":  info.FullMethod,
			"success": succ,
		},
	).Inc()
	return err
}

func unaryMiddle(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	requests_in_flight.
		With(prometheus.Labels{"method": info.FullMethod}).
		Inc()
	resp, err := handler(ctx, req)
	logging.LogRequest(info.FullMethod, err)
	succ := "yes"
	if err != nil {
		succ = "no"
	}
	requests_in_flight.
		With(prometheus.Labels{"method": info.FullMethod}).
		Dec()
	requests_total.With(
		prometheus.Labels{
			"method":  info.FullMethod,
			"success": succ,
		},
	).Inc()
	return resp, err
}
