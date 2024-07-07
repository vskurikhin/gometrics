/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * server_log.go
 * $Id$
 */

package interceptor

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/vskurikhin/gometrics/internal/logger"
)

func LogUnaryServer(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {

	logger.Log.Debug("LogUnaryServer", zap.String("request", fmt.Sprintf("%+v", req)))
	res, err := handler(ctx, req)
	logger.Log.Debug("LogUnaryServer", zap.String("response", fmt.Sprintf("%+v", res)))

	if err != nil {
		logger.Log.Debug("LogUnaryServer", zap.String("error", err.Error()))
	}
	return res, err
}
