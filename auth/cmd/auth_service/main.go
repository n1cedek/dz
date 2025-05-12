package main

import (
	"context"
	"dz/auth/app"
	"flag"
	"log"
)

const (
	idCreate = 1
)

//	func (s *UServer) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
//		log.Printf("Update user: %+v\n", req)
//		return &emptypb.Empty{}, nil
//	}
//
//	func (s *UServer) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
//		log.Printf("Delete user: %+v\n", req)
//		return &emptypb.Empty{}, nil
//	}
func main() {
	flag.Parse()
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}

}
