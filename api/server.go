package api

import (
	"fmt"
	db "simplebank/db/sqlc"
	"simplebank/util"

	"simplebank/token"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server servers HTTP request for our banking service
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine //this router will help to send each API req to the correct handler for processing
}

// Newserver creates a new HTTP server and setup routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	} // nil server

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// register the custom validator with gin (func IsSupportedCurrency)
	// convert the output to a validator.Validate pointer
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// if its okay then we call:
		// currency is the name of validation tag
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, err

}

func (server *Server) setupRouter() {
	router := gin.Default()

	// create new user
	router.POST("/users", server.createUser)
	// login request
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// create new account
	authRoutes.POST("/accounts", server.createAccount)
	// get aspecific account by ID
	authRoutes.GET("/accounts/:id", server.getAccount)
	// get list of account
	authRoutes.GET("/accounts", server.listAccount)
	// delete account by id
	authRoutes.DELETE("/accounts/:id", server.deleteAccount)
	// update account by id
	// router.PUT("/accounts/:id", server.updateAccount)
	// create a transfer
	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

// Start run the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func deleteResponse() gin.H {
	return gin.H{"succes": "account was deleted"}
}
