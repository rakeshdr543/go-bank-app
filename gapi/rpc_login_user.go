package gapi

import (
	"context"
	"time"

	db "github.com/rakeshdr543/go-bank-app/db/sqlc"
	pb "github.com/rakeshdr543/go-bank-app/pb/proto"
	"github.com/rakeshdr543/go-bank-app/util"
	"github.com/rakeshdr543/go-bank-app/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := ValidateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find user:%s", err)
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect password:%s", err)

	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(user.Username, util.DepositorRole, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to create token:%s", err)

	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(user.Username, util.DepositorRole, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to create refresh token:%s", err)

	}

	session, err := server.store.CreateSession(
		ctx, db.CreateSessionParams{
			ID:           refreshTokenPayload.ID,
			Username:     user.Username,
			RefreshToken: refreshToken,
			ClientIp:     server.extractMetadata(ctx).ClientIp,
			UserAgent:    server.extractMetadata(ctx).UserAgent,
			ExpiresAt:    time.Now().Add(server.config.RefreshTokenDuration),
		},
	)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to create session:%s", err)

	}

	resp := pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		User:                  convertUser(user),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiresAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiresAt),
	}

	return &resp, nil
}

func ValidateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUserName(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations

}
