package app

import (
	"context"
	"dz/auth/internal/closer"
	env "dz/auth/internal/config"
	"dz/auth/internal/interceptor"
	"dz/auth/internal/metric"
	descAccess "dz/auth/pkg/access_v1"
	desc "dz/auth/pkg/w1"
	"flag"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"sync"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to run GRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	go func() {
		err := runPrometheus()
		if err != nil {
			log.Fatal(err)
		}
	}()
	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx2 context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initMetrics,
		a.initGRPCServer,
		a.initHTTPServer,
	}
	for _, i := range inits {
		err := i(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func (a *App) initConfig(_ context.Context) error {
	err := env.Load(configPath)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initMetrics(ctx context.Context) error {
	return metric.Init(ctx)
}

func (a *App) initGRPCServer(ctx context.Context) error {
	cred, err := a.serviceProvider.TLSConfig().GetTLSConfig()
	if err != nil {
		log.Fatalf("failed to get TLS credentials: %v", err)
		return err
	}

	a.grpcServer = grpc.NewServer(
		grpc.Creds(cred),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			interceptor.ValidateInterceptor,
			interceptor.MetricsInterceptor,
		)))

	reflection.Register(a.grpcServer)
	desc.RegisterUserAPIServer(a.grpcServer, a.serviceProvider.AuthImpl(ctx))
	descAccess.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessI(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	creds, err := credentials.NewClientTLSFromFile("/home/n1cedek/GolandProjects/dz/auth/certificates/ca.cert", "")
	if err != nil {
		return err
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	err = desc.RegisterUserAPIHandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: mux,
	}
	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address())
	lis, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}
	err = a.grpcServer.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}
func (a *App) runHTTPServer() error {
	log.Printf("HTTP server is running on %s", a.serviceProvider.HTTPConfig().Address())
	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func runPrometheus() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	prometheusServer := &http.Server{
		Addr:    "localhost:2112",
		Handler: mux,
	}
	log.Printf("Prometheus server is running on %s", "localhost:2112")
	err := prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
