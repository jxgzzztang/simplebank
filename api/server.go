package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/jxgzzztang/simplebank/db/sqlc"
	"github.com/jxgzzztang/simplebank/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


type Server struct {
	store db.Store
	router *gin.Engine
}

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/

//	@securityDefinitions.basic	BasicAuth

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

func NewServer(store db.Store) Server {
	server := Server{
		store: store,
	}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", currencyValidate)
		if err != nil {
			return Server{}
		}
	}
	docs.SwaggerInfo.Title = "Simplebank API"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/login", server.Login)
	router.POST("/createUser", server.CreateUser)
	router.POST("/renewAccessToken", server.RenewAccessToken)
	RouterGroup(router, server)
	server.router = router
	return server
}

func RouterGroup(router *gin.Engine, server Server) {
	routerGroup := router.Group("/")
	routerGroup.Use(authMiddleware())
	routerGroup.GET("/account/:id", server.GetAccount)
	router.POST("/createAccount", server.CreateAccount)
	routerGroup.GET("/listAccounts", server.ListAccounts)
	routerGroup.POST("/transfer", server.Transfer)
} 

func (server *Server) Start(address string) error {
	gin.SetMode(gin.DebugMode)
	err := server.router.Run(address)
	return err
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func errorResponse(error error) ErrorResponse {
	return ErrorResponse{
		Error: error.Error(),
	}
}