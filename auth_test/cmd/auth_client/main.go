package main

import (
	"context"
	desc "dz/auth_test/pkg/w1"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to serve: %v", err)
	}
	defer conn.Close()

	c := desc.NewUserAPIClient(conn)

	createdID := makeCreate(ctx, c)

	makeGet(ctx, c, createdID)
	makeUpdate(ctx, c, createdID)
	makeGet(ctx, c, createdID)
	makeDelete(ctx, c, createdID)
}

func makeCreate(ctx context.Context, c desc.UserAPIClient) int64 {
	req := &desc.CreateRequest{
		Name:            gofakeit.BeerName(),
		Email:           gofakeit.Email(),
		Password:        "123",
		PasswordConfirm: "123",
		Role:            desc.Role_user,
	}

	resp, err := c.Create(ctx, req)
	if err != nil {
		log.Fatalf("failed to create user: %v\n", err)
	}
	log.Println(color.YellowString("Created user with ID: %d", resp.Id))
	return resp.Id
}

func makeUpdate(ctx context.Context, c desc.UserAPIClient, id int64) {
	req := &desc.UpdateRequest{
		Id:    id,
		Name:  wrapperspb.String(gofakeit.Name()),
		Email: wrapperspb.String(gofakeit.Email()),
	}

	_, err := c.Update(ctx, req)
	if err != nil {
		log.Fatalf("failed to update: %v", err)
	}
	log.Println(color.YellowString("Updated user with ID: %d", id))
}

func makeDelete(ctx context.Context, c desc.UserAPIClient, id int64) {
	req := &desc.DeleteRequest{Id: id}

	_, err := c.Delete(ctx, req)
	if err != nil {
		log.Fatalf("failed to delete user: %v", err)
	}

	log.Println(color.YellowString("Deleted user with ID: %d", id))
}

func makeGet(ctx context.Context, c desc.UserAPIClient, id int64) {
	req := &desc.GetRequest{Id: id}

	resp, err := c.Get(ctx, req)
	if err != nil {
		log.Fatalf("failed to get user: %v", err)
	}
	log.Printf(color.RedString("User info:\n", color.GreenString("%+v", resp)))
}
