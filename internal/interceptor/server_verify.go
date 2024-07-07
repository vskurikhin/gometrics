/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * server_verify.go
 * $Id$
 */

package interceptor

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/ip"
)

type VerifyXRealIP struct {
	IpNet *net.IPNet
}

func GetXRealIPVerifyer(cfg env.Config) func(
	context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler,
) (interface{}, error) {
	return (&VerifyXRealIP{
		IpNet: ip.TrustedIpNet(cfg),
	}).VerifyUnaryServer
}

func (v *VerifyXRealIP) VerifyUnaryServer(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {

	if v.IpNet != nil {

		md, ok := metadata.FromIncomingContext(ctx)

		if !ok || len(md["x-real-ip"]) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
		}
		if !v.IpNet.Contains(net.ParseIP(md["x-real-ip"][0])) {
			return nil, status.Error(codes.PermissionDenied, "unexpected client")
		}
	}
	res, err := handler(ctx, req)

	return res, err
}
