package consul

import (
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

type GRPCHealthCheckFunc func(ctx context.Context, client *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error)

type grpcHealth struct {
	check GRPCHealthCheckFunc
}

func (g *grpcHealth) Check(ctx context.Context, client *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	if g.check != nil {
		return g.check(ctx, client)
	}
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (*grpcHealth) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return nil
}
