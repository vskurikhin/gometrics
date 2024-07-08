/*
 * This file was last modified at 2024-07-08 14:51 by Victor N. Skurikhin.
 * grpc_reports.go
 * $Id$
 */

package agent

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"

	"github.com/vskurikhin/gometrics/internal/crypto"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/interceptor"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/types"

	pb "github.com/vskurikhin/gometrics/proto"
)

func grpcReports(cfg env.Config, enabled []types.Name) {

	request := new(pb.MetricsUpdateRequest)
	request.Metrics = new(pb.Metrics)
	request.Metrics.Metrics = make([]*pb.Metric, 0)

	for _, i := range enabled {
		metric := getMetric(i)
		if metric != nil {
			request.Metrics.Metrics = append(request.Metrics.Metrics, metric.ToResponse())
		}
	}
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		grpc.WithChainUnaryInterceptor(
			interceptor.XRealIP{Ip: cfg.Property().OutboundIP().String()}.UnaryClient,
		),
	}
	tlsCredentials, err := crypto.LoadAgentTLSCredentials()
	if err != nil {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		logger.Log.Debug("gRPC client load TLS credentials", zap.String("error", err.Error()))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	}
	conn, err := grpc.NewClient(cfg.GRPCAddress(), opts...)
	if err != nil {
		logger.Log.Debug("gRPC client connect", zap.String("error", err.Error()))
		return
	}
	defer func() { _ = conn.Close() }()
	c := pb.NewMetricsServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer func() {
		cancel()
		ctx.Done()
	}()
	resp, err := c.Updates(context.Background(), request)
	for i := 0; err != nil && isUpperBound(i, cfg.ReportInterval()); i++ {
		logger.Log.Debug("gRPC updates",
			zap.String("error", fmt.Sprintf("%v", err)),
			zap.String("time", fmt.Sprintf("%v", time.Now())),
		)
		resp, err = c.Updates(context.Background(), request)
	}
	if err == nil {
		logger.Log.Debug("gRPC client updates", zap.String("metrics", fmt.Sprintf("%s", request.Metrics.Metrics)))
		logger.Log.Debug("gRPC updates", zap.String("status", fmt.Sprintf("%v", resp)))
	}
}
