package svc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "{{ cookiecutter.go_module_path.strip('/') }}/api/go/{{ cookiecutter.project_name.strip() }}/v1"
	"{{ cookiecutter.go_module_path.strip('/') }}/internal/convert"
	"{{ cookiecutter.go_module_path.strip('/') }}/internal/repos/users"
)

func (s *Service) AddUser(
	ctx context.Context,
	req *pb.AddUserRequest,
) (*pb.AddUserResponse, error) {
	// Validate input
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}
	if req.Age <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "age must be positive")
	}

	now := time.Now()

	user, err := s.usersRepo.AddUser(ctx, users.AddUserParams{
		Name:      req.Name,
		Age:       req.Age,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add user: %v", err)
	}

	return &pb.AddUserResponse{
		User: convert.UserDBToProto(user),
	}, nil
}

func (s *Service) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (*pb.GetUserResponse, error) {
	user, err := s.usersRepo.GetUserById(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &pb.GetUserResponse{
		User: convert.UserDBToProto(user),
	}, nil
}

func (s *Service) UpdateUser(
	ctx context.Context,
	req *pb.UpdateUserRequest,
) (*pb.UpdateUserResponse, error) {
	// Validate input
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}
	if req.Age <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "age must be positive")
	}

	user, err := s.usersRepo.UpdateUserById(ctx, users.UpdateUserByIdParams{
		ID:        req.Id,
		Name:      req.Name,
		Age:       req.Age,
		UpdatedAt: time.Now(),
	}, &req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &pb.UpdateUserResponse{
		User: convert.UserDBToProto(user),
	}, nil
}

func (s *Service) DeleteUser(
	ctx context.Context,
	req *pb.DeleteUserRequest,
) (*pb.DeleteUserResponse, error) {
	user, err := s.usersRepo.DeleteUserById(ctx, req.Id, &req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &pb.DeleteUserResponse{
		User: convert.UserDBToProto(user),
	}, nil
}
