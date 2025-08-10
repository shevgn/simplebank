// Package api provides API endpoints for the simplebank application.
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/shevgn/simplebank/db/sqlc"
	"github.com/shevgn/simplebank/token"
	"github.com/shevgn/simplebank/util"
)

// Server represents the API server.
type Server struct {
	store      db.Store
	tokenMaker token.Maker
	config     *util.Config
	router     *gin.Engine
}

// NewServer creates a new API server.
func NewServer(config *util.Config, store db.Store) *Server {
	maker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		panic(err)
	}

	s := &Server{
		store:      store,
		tokenMaker: maker,
		config:     config,
		router:     gin.Default(),
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			panic(err)
		}
	}

	s.registerRoutes()

	return s
}

func (s *Server) registerRoutes() {
	s.router.POST("/users", s.createUser)
	s.router.POST("/users/login", s.loginUser)
	s.router.POST("/tokens/renew_access", s.renewAccessToken)

	authRoutes := s.router.Group("/").Use(authMiddleware(s.tokenMaker))

	authRoutes.GET("/accounts/:id", s.getAccount)
	authRoutes.GET("/accounts", s.listAccounts)
	authRoutes.POST("/accounts", s.createAccount)
	authRoutes.PUT("/accounts", s.updateAccount)
	authRoutes.DELETE("/accounts/:id", s.deleteAccount)

	authRoutes.POST("/transfers", s.createTransfer)
}

// Run starts the API server.
func (s *Server) Run(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
