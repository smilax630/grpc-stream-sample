package grpc

import (
	"context"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/grpc-streamer/pkg/util/metrics"
)

const name = "grpc"

var (
	registerOnce sync.Once
	latencyMs    = stats.Int64(name, "The grpc latency in milliseconds", stats.UnitMilliseconds)
	counter      = stats.Int64(name, "The grpc call count", stats.UnitDimensionless)
)

var (
	keyRPCType = tag.MustNewKey("rpc_type")
	keyService = tag.MustNewKey("service")
	keyMethod  = tag.MustNewKey("method")
	keyStatus  = tag.MustNewKey("status")
)

func RegisterMetrics() error {
	latencyView := &view.View{
		Name:        "grpc_latency_distribution",
		Measure:     latencyMs,
		Description: "The distribution of grpc API latencies",
		TagKeys:     []tag.Key{keyRPCType, keyService, keyMethod, keyStatus},
		Aggregation: metrics.NewLatencyDistribution(),
	}
	countView := &view.View{
		Name:        "grpc_api_count",
		Measure:     counter,
		Description: "The count of grpc API calls",
		TagKeys:     []tag.Key{keyRPCType, keyService, keyMethod, keyStatus},
		Aggregation: view.Count(),
	}
	var err error
	registerOnce.Do(func() {
		err = view.Register(latencyView, countView)
	})
	return err
}

func NewUnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	const rpcType = "unary"
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		start := time.Now()
		service, method := splitFullMethodName(info.FullMethod)
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Internal, "%v, stack:%s", r, debug.Stack())
			}

			code := status.Code(err).String()
			latency := time.Since(start).Nanoseconds() / 1e6 // nolint: gomnd
			taggedCtx, err := newContext(rpcType, service, method, code)
			if err != nil {
				logger.Error("Failed to create context", zap.Error(err))
			}
			count(taggedCtx)
			measure(taggedCtx, latency)
		}()
		return handler(ctx, req)
	}
}

func splitFullMethodName(fullMethodName string) (string, string) {
	// format: /project.package.service/method
	const partsCount = 3
	parts := strings.Split(fullMethodName, "/")
	if len(parts) != partsCount {
		return "unknown", "unknown"
	}
	return parts[1], parts[2]
}

func newContext(rpcType, service, method, status string) (context.Context, error) {
	return tag.New(
		context.Background(),
		tag.Insert(keyRPCType, rpcType),
		tag.Insert(keyService, service),
		tag.Insert(keyMethod, method),
		tag.Insert(keyStatus, status),
	)
}

func count(ctx context.Context) {
	stats.Record(ctx, counter.M(1))
}

func measure(ctx context.Context, ms int64) {
	stats.Record(ctx, latencyMs.M(ms))
}
