package gapi

import (
	db "github.com/rakeshdr543/go-bank-app/db/sqlc"
	pb "github.com/rakeshdr543/go-bank-app/pb/proto"
	"github.com/rakeshdr543/go-bank-app/token"
	"github.com/rakeshdr543/go-bank-app/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      *db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store *db.Store) (*Server, error) {

	tokenMaker, err := token.NewJWTMaker(
		config.TokenSymmetricKey,
	)

	if err != nil {
		return nil, err
	}

	server := &Server{store: store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	return server, nil
}
