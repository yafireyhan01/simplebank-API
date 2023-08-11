package api

import (
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server servers HTTP request for our banking service
type Server struct {
	store  db.Store
	router *gin.Engine //this router will help to send each API req to the correct handler for processing
}

// Newserver creates a new HTTP server and setup routing.
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// register the custom validator with gin (func IsSupportedCurrency)
	// convert the output to a validator.Validate pointer
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// if its okay then we call:
		// currency is the name of validation tag
		v.RegisterValidation("currency", validCurrency)
	}

	// create new account
	router.POST("/accounts", server.createAccount)
	// get aspecific account by ID
	router.GET("/accounts/:id", server.getAccount)
	// get list of account
	router.GET("/accounts", server.listAccount)
	// delete account by id
	router.DELETE("/accounts/:id", server.deleteAccount)
	// update account by id
	// router.PUT("/accounts/:id", server.updateAccount)
	// create a transfer
	router.POST("/transfers", server.createTransfer)
	// create new user
	router.POST("/users", server.createUser)

	server.router = router
	return server

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
