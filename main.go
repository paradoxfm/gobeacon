package main

import (
	"gobeacon/controller"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "gobeacon/docs"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
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

	srvSwag := createSwaggerApi()
	srvPhone := createPhoneApi()
	srvPhoneAdm := createPhoneAdminApi()

	for _, value := range []*http.Server{srvSwag, srvPhone, srvPhoneAdm} {
		g.Go(func() error {
			return value.ListenAndServe()
		})
	}
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func createPhoneApi() (*http.Server) {
	r := gin.New()
	auth := controller.CreateHeartGinJWTMiddleware()
	r.Use(auth.MiddlewareFunc())
	r.GET("/api/v1/heartbeat", dummyHandler)
	return initServer(":7777", r)
}

func createPhoneAdminApi() (*http.Server) {
	auth := controller.CreateAdminJWTMiddleware()
	r := gin.New()
	v1 := r.Group("/api/v1")  // api первой версии
	usr := v1.Group("/users") // api для пользователей
	mFunc := auth.MiddlewareFunc()
	{
		usr.POST("/signUp", controller.UserCreate)
		usr.POST("/login", auth.LoginHandler)
		usr.POST("/password/reset", dummyHandler)
		me := usr.Group("/me")
		me.Use(mFunc)
		{
			usr.GET("", dummyHandler)
			usr.PUT("/password", dummyHandler)
			usr.PUT("/push", dummyHandler)
			usr.PUT("/refresh", auth.RefreshHandler)
		}
	}

	trk := v1.Group("/trackers") // api для трекеров
	trk.Use(mFunc)
	{
		trk.POST("", dummyHandler)
		trk.GET("/:id", dummyHandler)
		trk.DELETE("/:id", dummyHandler)
		trk.PUT("/:id", dummyHandler)
		trk.POST("/:id/avatar", dummyHandler)
		trk.GET("/:id/geo", dummyHandler)
		trk.GET("/:id/geo/history", dummyHandler) //date_start date_end
	}

	zone := v1.Group("/geoZones") // api для гео зон
	zone.Use(mFunc)
	{
		zone.GET("", dummyHandler)
		zone.POST("", dummyHandler)
		zone.DELETE("/:id", dummyHandler)
		zone.GET("/:id", dummyHandler)
		zone.PUT("/:id", dummyHandler)
		zone.PUT("/:id/trackers", dummyHandler)
	}

	return initServer(":8070", r)
}

func createWatchApi() (*http.Server) {
	r := gin.New()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return initServer(":6666", r)
}

func createSwaggerApi() (*http.Server) {
	r := gin.New()
	// документация по сервисам /swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return initServer(":8071", r)
}

func dummyHandler(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}

func initServer(port string, routes *gin.Engine) (*http.Server) {
	srv := &http.Server{
		Addr:    port,
		Handler: routes,
	}
	return srv
}
