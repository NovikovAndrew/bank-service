package api

import (
	"bank-service/db/sqlc"
	"bank-service/token"
	"bank-service/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves all HTTP request for banking service
type Server struct {
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
	router     *gin.Engine
}

func NewServer(store db.Store, tokenMaker token.Maker, config util.Config) *Server {
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("currency", currencyValidator)
	}

	// add routes to router
	server.setupRouter()

	return server
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// users
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/token/renew_access", server.renewAccessToken)

	authGroup := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// accounts
	authGroup.POST("/accounts", server.createAccount)
	authGroup.GET("/accounts/:id", server.getAccount)
	authGroup.GET("/accounts", server.listAccount)
	authGroup.PUT("/accounts", server.updateAccount)
	authGroup.DELETE("/accounts/:id", server.deleteAccount)

	// transfers
	authGroup.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
