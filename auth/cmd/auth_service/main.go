package main

import (
	"context"
	config "dz/auth/internal"
	env "dz/auth/internal/config"
	desc "dz/auth/pkg/w1"
	"flag"
	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

const (
	idCreate = 1
)

var configPath string

type UServer struct {
	desc.UnimplementedUserAPIServer
}

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func (s *UServer) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("ID: %d", req.GetId())
	return &desc.GetResponse{
		Id:        req.GetId(),
		Name:      gofakeit.BeerName(),
		Email:     gofakeit.Email(),
		Role:      desc.Role_user,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

func (s *UServer) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Created user: %+v\n", req)
	return &desc.CreateResponse{Id: idCreate}, nil
}

func (s *UServer) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("Update user: %+v\n", req)
	return &emptypb.Empty{}, nil
}

func (s *UServer) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Delete user: %+v\n", req)
	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	appConfig, err := env.NewAppConfig()
	if err != nil {
		log.Fatalf("failed to get db config: %v", err)
	}

	//connect to database
	con, err := pgx.Connect(ctx, appConfig.DB.DSN())
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer con.Close(ctx)

	//request insert
	builderInsert := sq.Insert("auth").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email").
		Values(gofakeit.BeerName(), gofakeit.Email()).
		Suffix("RETURNING id")

	query, arg, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to builde query: %v", err)
	}

	var noteID int
	err = con.QueryRow(ctx, query, arg...).Scan(&noteID)
	if err != nil {
		log.Fatalf("failed to insert note: %v", err)
	}
	log.Printf("inserted note with id: %d", noteID)

	//lis, err := net.Listen("tcp", fmt.Sprintf(":%d", appConfig.CC.Port()))
	//if err != nil {
	//	log.Fatalf("failed to listen server: %v", err)
	//}
	//s := grpc.NewServer()
	//reflection.Register(s)
	//desc.RegisterUserAPIServer(s, &UServer{})
	//
	//log.Printf("Server listening at: %v", lis.Addr())
	//
	//if err = s.Serve(lis); err != nil {
	//	log.Fatalf("failed to run server: %v", err)
	//}

}
