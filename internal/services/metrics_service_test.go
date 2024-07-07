/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * metrics_service_test.go
 * $Id$
 */

package services

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/credentials/local"

	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/types"

	"google.golang.org/grpc/credentials/insecure"

	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/vskurikhin/gometrics/internal/util"

	pb "github.com/vskurikhin/gometrics/proto"
)

var (
	testDataBaseDSN   = ""
	testKey           string
	testPort          int
	testServerAddress string
	testTempFileName  string
)

func TestMetricsService(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cfg := getTestConfig()
	go grpcServe(ctx, cfg)
	conn := grpcClientConn()
	defer func() { _ = conn.Close() }()

	tests := []struct {
		name string
		fTst func(t *testing.T, c pb.MetricsServiceClient)
	}{
		{
			name: "positive test #1 update",
			fTst: MetricsServiceClientUpdateTest,
		},
		{
			name: "positive test #2 updates",
			fTst: MetricsServiceClientUpdatesTest,
		},
		{
			name: "positive test #3 value",
			fTst: MetricsServiceClientValueTest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fTst(t, pb.NewMetricsServiceClient(conn))
		})
	}
}

func MetricsServiceClientUpdateTest(t *testing.T, c pb.MetricsServiceClient) {
	metrics := []*pb.Metric{
		{Id: "Alloc", Type: "gauge", Value: 2.718},
		{Id: "PollCount", Type: "counter", Delta: 2},
	}
	for _, metric := range metrics {
		resp, err := c.Update(context.Background(), &pb.MetricUpdateRequest{
			Metric: metric,
		})
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		if resp.Status != pb.Status_OK {
			t.Fatal(resp.Error)
		}
	}
}

func MetricsServiceClientUpdatesTest(t *testing.T, c pb.MetricsServiceClient) {

	request := new(pb.MetricsUpdateRequest)
	request.Metrics = new(pb.Metrics)
	request.Metrics.Metrics = []*pb.Metric{
		{Id: "Alloc", Type: "gauge", Value: 3.14},
		{Id: "PollCount", Type: "counter", Delta: 3},
	}
	resp, err := c.Updates(context.Background(), request)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if resp.Status != pb.Status_OK {
		t.Fatal(resp.Error)
	}
}

func MetricsServiceClientValueTest(t *testing.T, c pb.MetricsServiceClient) {
	metrics := []struct {
		input *pb.MetricRequestValue
		want  interface{}
	}{
		{
			&pb.MetricRequestValue{Id: "Alloc", Type: "gauge"},
			3.14,
		},
		{
			&pb.MetricRequestValue{Id: "PollCount", Type: "counter"},
			int64(5),
		},
	}
	for _, metric := range metrics {
		resp, err := c.Value(context.Background(), &pb.MetricValueRequest{
			Metric: metric.input,
		})
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		if resp.Status != pb.Status_OK {
			t.Fatal(resp.Error)
		}
		switch {
		case types.GAUGE.Eq(resp.Metric.Type):
			assert.Equal(t, metric.want, resp.Metric.Value)
		case types.COUNTER.Eq(resp.Metric.Type):
			assert.Equal(t, metric.want, resp.Metric.Delta)
		}
	}
}

func grpcClientConn() *grpc.ClientConn {
	conn, err := grpc.NewClient(
		fmt.Sprintf("127.0.0.1:%d", testPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	util.IfErrorThenPanic(err)
	return conn
}

func grpcServe(ctx context.Context, cfg env.Config) {
	listen, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", testPort))
	util.IfErrorThenPanic(err)
	sopts := []grpc.ServerOption{grpc.Creds(local.NewCredentials())}
	s := grpc.NewServer(sopts...)
	ms := GetMetricsService(cfg)
	pb.RegisterMetricsServiceServer(s, ms)
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.Stop()
				return
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

func getTestConfig() env.Config {
	return env.GetTestConfig(
		env.GetProperty,
		env.WithDataBaseDSN(&testDataBaseDSN),
		env.WithFileStoragePath(testTempFileName),
		env.WithKey(&testKey),
		env.WithPollInterval(30*time.Minute),
		env.WithReportInterval(time.Hour),
		env.WithRestore(true),
		env.WithServerAddress(testServerAddress),
		env.WithStoreInterval(24*time.Hour),
	)
}

func init() {
	testPort = 65500 + rand.Intn(34)
	testKey = fmt.Sprintf("%018d", rand.Uint32())
	testServerAddress = fmt.Sprintf("localhost:%d", testPort)
	testTempFileName = fmt.Sprintf("%s/test_%018d.txt", os.TempDir(), rand.Uint32())
}
