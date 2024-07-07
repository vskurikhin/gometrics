/*
 * This file was last modified at 2024-07-08 13:55 by Victor N. Skurikhin.
 * metrics_service.go
 * $Id$
 */

// Package services реализация gRPC сервера.
package services

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vskurikhin/gometrics/internal/dto"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/storage"
	"github.com/vskurikhin/gometrics/internal/types"

	pb "github.com/vskurikhin/gometrics/proto"
)

type MetricsService interface {
	pb.MetricsServiceServer
	DTOUpdate(context.Context, *dto.Metric) (*dto.Metric, error)
	DTOUpdates(context.Context, *dto.Metrics) (*dto.Metrics, error)
	DTOValue(context.Context, *dto.Metric) (*dto.Metric, error)
}

type metricsService struct {
	pb.UnimplementedMetricsServiceServer
	store storage.Storage
}

var _ MetricsService = (*metricsService)(nil)
var onceMetricsService = new(sync.Once)
var metricsSrv *metricsService

// GetMetricsService — MetricsService сервис по работе с метриками,
// реализован как паттерн Front Controller (Контроллер запросов).
func GetMetricsService(cfg env.Config) MetricsService {

	onceMetricsService.Do(func() {
		metricsSrv = new(metricsService)
		metricsSrv.store = cfg.Property().Storage()
	})
	return metricsSrv
}

// DTOUpdate реализует интерфейс для использования в хендлере HTTP сервера по получению метрики из БД.
func (s *metricsService) DTOUpdate(ctx context.Context, metric *dto.Metric) (*dto.Metric, error) {
	return s.updateMetric(ctx, metric)
}

// DTOUpdates реализует интерфейс для использования в хендлере HTTP сервера по обновлению метрик в БД.
func (s *metricsService) DTOUpdates(ctx context.Context, metrics *dto.Metrics) (*dto.Metrics, error) {
	s.updateMetrics(ctx, *metrics)
	return metrics, nil
}

// DTOValue реализует интерфейс для использования в хендлере HTTP сервера по получению метрики из БД.
func (s *metricsService) DTOValue(ctx context.Context, metric *dto.Metric) (*dto.Metric, error) {
	return s.valueMetric(ctx, metric)
}

// Update реализует интерфейс gRPC сервера по обновлению метрики в БД.
func (s *metricsService) Update(ctx context.Context, in *pb.MetricUpdateRequest) (*pb.MetricUpdateResponse, error) {

	metric := dto.FromRequest(in.Metric)
	metric, err := s.updateMetric(ctx, metric)

	if err != nil {
		response := new(pb.MetricUpdateResponse)
		response.Status = pb.Status_FAIL
		response.Error = err.Error()
		return response, status.Errorf(codes.Aborted, err.Error())
	}
	response := new(pb.MetricUpdateResponse)
	response.Metric = metric.ToResponse()
	response.Status = pb.Status_OK

	return response, nil
}

// Updates реализует интерфейс gRPC сервера по обновлению метрик в БД.
func (s *metricsService) Updates(ctx context.Context, in *pb.MetricsUpdateRequest) (*pb.MetricsUpdateResponse, error) {

	metrics := make(dto.Metrics, 0)

	for _, metric := range in.GetMetrics().Metrics {
		metrics = append(metrics, *dto.FromRequest(metric))
	}
	s.updateMetrics(ctx, metrics)
	logger.Log.Debug("gRPC server updates", zap.String("status", metrics.String()))

	return &pb.MetricsUpdateResponse{
		Status: pb.Status_OK,
	}, nil
}

// Value реализует интерфейс gRPC сервера по получению метрики из БД.
func (s *metricsService) Value(ctx context.Context, in *pb.MetricValueRequest) (*pb.MetricValueResponse, error) {

	metric := dto.FromValueRequest(in.Metric)
	metric, err := s.valueMetric(ctx, metric)

	if err != nil {
		response := new(pb.MetricValueResponse)
		response.Status = pb.Status_FAIL
		response.Error = err.Error()
		return response, status.Errorf(codes.Aborted, err.Error())
	}
	response := new(pb.MetricValueResponse)
	response.Metric = metric.ToResponse()
	response.Status = pb.Status_OK

	return response, nil
}

func (s *metricsService) updateMetric(_ context.Context, metric *dto.Metric) (*dto.Metric, error) {

	var err error
	var name string
	num := types.Lookup(metric.ID)

	if num > 0 {
		name = num.String()
	} else {
		name = metric.ID
	}
	switch {
	case types.GAUGE.Eq(metric.MType):
		value := fmt.Sprintf("%.12f", *metric.Value)
		s.store.PutGauge(name, &value)
	case types.COUNTER.Eq(metric.MType):
		pv := s.store.GetCounter(name)
		*metric.Delta = metric.CalcDelta(pv)
		value := fmt.Sprintf("%d", *metric.Delta)
		s.store.PutCounter(name, &value)
	default:
		err = fmt.Errorf("update %s for type %s not found", name, metric.MType)
	}
	return metric, err
}

func (s *metricsService) updateMetrics(_ context.Context, metrics dto.Metrics) {
	s.store.PutSlice(metrics)
}

func (s *metricsService) valueMetric(_ context.Context, metric *dto.Metric) (*dto.Metric, error) {

	var err error
	var name string
	num := types.Lookup(metric.ID)

	if num > 0 {
		name = num.String()
	} else {
		name = metric.ID
	}
	switch {
	case types.GAUGE.Eq(metric.MType):
		value := s.store.GetGauge(name)
		metric.Value = new(float64)
		if value != nil {
			*metric.Value, err = strconv.ParseFloat(*value, 64)
		} else {
			err = fmt.Errorf("value %s of type %s not found", name, metric.MType)
		}
	case types.COUNTER.Eq(metric.MType):
		value := s.store.GetCounter(name)
		metric.Delta = new(int64)
		if value != nil {
			*metric.Delta, err = strconv.ParseInt(*value, 10, 64)
		} else {
			err = fmt.Errorf("value %s of type %s not found", name, metric.MType)
		}
	}
	return metric, err
}
