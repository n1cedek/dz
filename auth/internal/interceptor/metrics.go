package interceptor

import (
	"context"
	"dz/auth/internal/metric"
	"google.golang.org/grpc"
	"time"
)

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	metric.IncRequestCounter()

	timeStart := time.Now()
	res, err := handler(ctx, req)
	diffTime := time.Since(timeStart)
	if err != nil {
		metric.IncResponseCounter("error", info.FullMethod)
		metric.IncHistogramCounter("error", diffTime.Seconds())
	} else {
		metric.IncResponseCounter("success", info.FullMethod)
		metric.IncHistogramCounter("success", diffTime.Seconds())
	}
	return res, err
}
