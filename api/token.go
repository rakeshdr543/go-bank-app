package api

import (
	"net/http"

	"github.com/rakeshdr543/go-bank-app/token"
	"github.com/rakeshdr543/go-bank-app/util"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type renewTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type renewTokenResponse struct {
	SessionId            uuid.UUID    `json:"session_id"`
	User                 userResponse `json:"user"`
	AccessToken          string       `json:"access_token"`
	AccessTokenExpiresAt time.Time    `json:"access_token_expires_at"`
}

func (server *Server) renewToken(ctx *gin.Context) {
	var req renewTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := token.Maker.VerifyToken(server.tokenMaker, req.RefreshToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, payload.ID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, session.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(user.Username, util.DepositorRole, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := renewTokenResponse{
		SessionId:            session.ID,
		User:                 newUserResponse(user),
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiresAt,
	}

	ctx.JSON(http.StatusOK, resp)
}
