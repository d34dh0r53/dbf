package main

import (
	"context"
	"log"
	"net"

	pb "github.com/d34dh0r53/dbf/dbf-service/proto/osa-overrides"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":42424"
)

// OverridesInterface - interface for overrides
type OverridesInterface interface {
	Fetch(*pb.OSASha) (*pb.DefaultsFile, error)
}

// DefaultsRepo - placeholder for defaults
type DefaultsRepo struct {
	defaults []*pb.DefaultsFile
}

type service struct {
	overrides OverridesInterface
}

// Fetch - method to fetch defaults from a given OSA SHA
func (repo *DefaultsRepo) Fetch(osasha *pb.OSASha) (*pb.DefaultsFile, error) {
	var test *pb.DefaultsFile
	test.Path = "/var/test"
	test.Contents = "This is a test"
	return test, nil
}

func (s *service) GetOverrides(ctx context.Context, sha *pb.OSASha) (*pb.OSADefaults, error) {
	// do something with the DefaultsFile
	overrides, err := s.overrides.Fetch(sha)
	if err != nil {
		return nil, err
	}

	return &pb.OSADefaults{Valid: true, Defaultsfile: overrides}, nil
}

func main() {

	repo := &DefaultsRepo{}

	// Set up our gRPC Server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterOverrideGeneratorServer(s, &service{repo})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
