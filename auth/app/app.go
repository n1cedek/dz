package app

import (
	"context"
	"dz/auth/closer"
	env "dz/auth/internal/config"
	desc "dz/auth/pkg/w1"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
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
	return a.RunGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx2 context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
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

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)
	desc.RegisterUserAPIServer(a.grpcServer, a.serviceProvider.AuthImpl(ctx))
	return nil
}

func (a *App) RunGRPCServer() error {
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
