/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * x_real_ip.go
 * $Id$
 */

package interceptor

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/vskurikhin/gometrics/internal/logger"
)

type XRealIP struct {
	Ip string
}

func (i XRealIP) UnaryClient(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {

	md := metadata.Pairs()
	md.Set("x-real-ip", i.Ip)
	ctx = metadata.NewOutgoingContext(ctx, md)
	err := invoker(ctx, method, req, reply, cc, opts...)

	if err != nil {
		logger.Log.Debug("XRealIP.UnaryClient", zap.String("error", err.Error()))
	}
	logger.Log.Debug("XRealIP.UnaryClient", zap.String("x-real-ip", i.Ip))

	return err
}
