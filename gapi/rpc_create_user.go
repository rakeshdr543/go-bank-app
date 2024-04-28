package gapi

import (
	"context"

	db "github.com/rakeshdr543/go-bank-app/db/sqlc"
	pb "github.com/rakeshdr543/go-bank-app/pb/proto"
	"github.com/rakeshdr543/go-bank-app/util"
	"github.com/rakeshdr543/go-bank-app/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := ValidateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password:%s", err)
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		FullName: req.FullName,
		Email:    req.Email,
		Password: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user:%s", err)
	}

	resp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return resp, nil

}

func ValidateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUserName(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateUserFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	return violations

}
