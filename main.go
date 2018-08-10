package main

import (
	"gobeacon/controller"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	_ "gobeacon/docs"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.basic BasicAuth
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	port := os.Getenv("PORT")
	r := gin.Default()

	if port == "" {
		port = "8000"
	}

	authMiddleware := controller.CreateGinJWTMiddleware()

	r.POST("/login", authMiddleware.LoginHandler)

	v1Users := r.Group("/api/v1/users")
	{
		v1Users.POST("/signup", controller.CreateUser)
	}

	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", controller.HelloHandler)
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
		auth.GET("/logout", authMiddleware.RefreshHandler)
	}
	// документация по сервисам /swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
