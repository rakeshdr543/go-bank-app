package api

import (
	db "sampla_bank/db/sqlc"
	"sampla_bank/token"
	"sampla_bank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store *db.Store) *Server {

	tokenMaker, err := token.NewJWTMaker(
		config.TokenSymmetricKey,
	)

	if err != nil {
		panic(err)
	}

	server := &Server{store: store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users/login", server.loginUser)
	router.POST("/users", server.createUser)

	groupRouter := router.Group("/", authMiddleware(server.tokenMaker))
	groupRouter.GET("/users/:username", server.getUser)

	groupRouter.POST("/accounts", server.createAccount)
	groupRouter.GET("/accounts/:id", server.getAccount)
	groupRouter.GET("/accounts", server.listAccounts)

	groupRouter.POST("/transfers", server.createMoneyTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
