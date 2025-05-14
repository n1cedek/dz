package app

import (
	"context"
	auth2 "dz/auth/internal/api/auth"
	"dz/auth/internal/closer"
	env "dz/auth/internal/config"
	"dz/auth/internal/repo"
	authRepo "dz/auth/internal/repo/auth"
	"dz/auth/internal/service"
	authServ "dz/auth/internal/service/auth"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type serviceProvider struct {
	pgConfig    env.PGConfig
	grpcConfig  env.GPRCConfig
	httpConfig  env.HTTPConfig
	tlsConfig   env.TLSConfig
	pgPool      *pgxpool.Pool
	authService service.AuthService
	authRepo    repo.AuthRepo
	authImpl    *auth2.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GRPCConfig() env.GPRCConfig {
	if s.grpcConfig == nil {
		grpcConfig, err := env.NewGrpcConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}
		s.grpcConfig = grpcConfig
	}
	return s.grpcConfig
}
func (s *serviceProvider) HTTPConfig() env.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHttpConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}
func (s *serviceProvider) PGConfig() env.PGConfig {
	if s.pgConfig == nil {
		pgConfig, err := env.NewDsnConfig()
		if err != nil {
			log.Fatalf("failed to get db config: %v", err)
		}
		s.pgConfig = pgConfig
	}
	return s.pgConfig
}

func (s *serviceProvider) TLSConfig() env.TLSConfig {
	if s.tlsConfig == nil {
		tls, err := env.NewTLSConfig()
		if err != nil {
			log.Fatalf("failed to initialize TLS config: %v", err)
		}
		s.tlsConfig = tls
	}
	return s.tlsConfig
}

func (s *serviceProvider) PgPool(ctx context.Context) *pgxpool.Pool {
	if s.pgPool == nil {
		con, err := pgxpool.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to db: %v", err)
		}
		err = con.Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %v", err)
		}

		closer.Add(func() error {
			con.Close()
			return nil
		})
		s.pgPool = con
	}
	return s.pgPool
}

func (s *serviceProvider) AuthRepo(ctx context.Context) repo.AuthRepo {
	if s.authRepo == nil {
		s.authRepo = authRepo.NewRepo(s.PgPool(ctx))
	}
	return s.authRepo
}
func (s *serviceProvider) AuthServ(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authServ.NewServ(s.AuthRepo(ctx))
	}
	return s.authService
}
func (s *serviceProvider) AuthImpl(ctx context.Context) *auth2.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth2.NewImplementation(s.AuthServ(ctx))
	}
	return s.authImpl
}
