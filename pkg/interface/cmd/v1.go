package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"

	"github.com/go-redis/redis"
	"github.com/grpc-streamer/config"
	"github.com/grpc-streamer/pkg/gateway"
	"github.com/grpc-streamer/pkg/usecase"
	pb "github.com/grpc-streamer/proto"
)

const port = ":50052"

// v1 means di container
type v1 struct {
	cmd               *cobra.Command
	redis             *redis.Client
	streamUseCase     *usecase.StreamUseCase
	callStatusGateway *gateway.CallStatusRedisGateway
}

//env is get environment ex) dev stg load prd
var env = os.Getenv("ENV")

func newV1() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "v1",
		Short: "Barry service server v1",
	}
	app := &v1{cmd: cmd}
	cmd.RunE = func(_ *cobra.Command, _ []string) error {
		return app.run()
	}
	return cmd
}

func (v *v1) run() error {
	// serverのtimezoneをutcに固定
	time.Local = time.FixedZone("UTC", 0)

	//serverの初期化
	sv := &v1{}

	fmt.Fprint(os.Stdout, "set up gRPC v1\n")

	config.SetupConfig()

	sv.setRedis()
	sv.redis.Ping()
	defer sv.redis.Close()

	//DI settings
	sv.setInjector()

	fmt.Fprint(os.Stdout, "start gRPC v1\n")

	sv.newGrpcServer()
	return nil
}

//setRedis is connection masterDB setting
func (v *v1) setRedis() {
	var err error
	redis := config.GetRedisConfig()
	v.redis, err = gateway.NewRedisClient(gateway.RedisOption{
		Hosts:       redis.Host,
		Database:    redis.DB,
		KeyPrefix:   redis.KeyPrefix,
		Password:    redis.Password,
		DialTimeout: 3 * time.Second,
		PoolSize:    redis.PoolSize,
	})
	if err != nil {
		log.Fatalf("NewRedisClient is failed. caz %v", err)
	}
}

// setInjector is solve the dependency injection
func (v *v1) setInjector() {
	v.callStatusGateway = gateway.NewCallStatusRedisGateway(v.redis)

	//use case
	v.streamUseCase = usecase.NewStreamUsecase(*v.callStatusGateway)
}

type httpHandler struct {
	grpcServer *grpc.Server
}

// newGrpcServer is start grpc v1
func (v *v1) newGrpcServer() {
	// gRPC new
	sv := grpc.NewServer()

	healthpb.RegisterHealthServer(sv, health.NewServer())

	pb.RegisterStreamServiceServer(sv, v.streamUseCase)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		if err := sv.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	<-sigCh
	sv.GracefulStop()
	fmt.Fprint(os.Stdout, "Shutting down...\n")

}

// UnaryServerInterceptor is auth checker.
func (v *v1) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = grpc.SendHeader(ctx, metadata.New(map[string]string{
			"Cache-control": "no-store, no-cache",
		}))

		resp, err = handler(ctx, req)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
