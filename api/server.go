// Package api provides API endpoints for the simplebank application.
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/shevgn/simplebank/db/sqlc"
)

// Server represents the API server.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new API server.
func NewServer(store db.Store) *Server {
	s := &Server{
		store:  store,
		router: gin.Default(),
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			panic(err)
		}
	}

	s.router.GET("/accounts/:id", s.getAccount)
	s.router.GET("/accounts", s.listAccounts)
	s.router.POST("/accounts", s.createAccount)
	s.router.PUT("/accounts", s.updateAccount)
	s.router.DELETE("/accounts/:id", s.deleteAccount)

	s.router.POST("/transfers", s.createTransfer)

	return s
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
