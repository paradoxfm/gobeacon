package main

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gobeacon/controller"
	_ "gobeacon/docs"
	"log"
	"net/http"
	"sync"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.basic BasicAuth
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	servers := map[string]http.Handler{":7777": createPhoneApi(), ":8070": createPhoneAdminApi(), ":8071": createSwaggerApi()}
	var wg sync.WaitGroup
	wg.Add(1)
	for port, server := range servers {
		go func(port string, server http.Handler) (err error) {
			defer log.Fatal(err)
			err = http.ListenAndServe(port, server)
			wg.Done()
			return err
		}(port, server)
	}
	wg.Wait()
}

func createSwaggerApi() (*gin.Engine) {
	r := gin.New()
	r.Use(gin.Recovery())
	// документация по сервисам /swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//r.Run(":8071")
	//http.ListenAndServe(":8071", r)
	return r //initServer(":8071", r)
}

func createPhoneApi() (*gin.Engine) {
	r := gin.New()
	r.Use(gin.Recovery())
	auth := controller.CreateHeartGinJWTMiddleware()
	r.Use(auth.MiddlewareFunc())
	r.GET("/api/v1/heartbeat", controller.HeartbeatPhone)
	return r // initServer(":7777", r)
}

func createPhoneAdminApi() (*gin.Engine) {
	auth := controller.CreateAdminJWTMiddleware()
	r := gin.Default()
	r.Use(gin.Recovery())
	v1 := r.Group("/api/v1")  // api первой версии
	usr := v1.Group("/users") // api для пользователей
	mFunc := auth.MiddlewareFunc()
	tst := v1.Group("/test")
	tst.Use(mFunc)
	tst.GET("/push", controller.TestPush)
	{
		usr.POST("/signUp", controller.UserCreate)
		usr.POST("/login", auth.LoginHandler)
		usr.POST("/password/reset", controller.UserResetPassword)
		me := usr.Group("/me")
		me.Use(mFunc)
		{
			me.GET("", controller.UserGetProfile)
			me.PUT("/password", controller.UserChangePassword)
			me.PUT("/push", controller.UserUpdatePushId)
			me.PUT("/avatar", controller.UserUpdateAvatar)
			me.PUT("/refresh", auth.RefreshHandler)
		}
	}

	trk := v1.Group("/trackers") // api для трекеров
	trk.Use(mFunc)
	{
		trk.GET("/all", controller.TrackAll)
		trk.POST("/custom", controller.TrackByIds)
		trk.POST("", controller.TrackCreate)
		trk.GET("/find/:id", controller.TrackGetById)
		trk.DELETE("/delete/:id", controller.TrackDeleteById)
		trk.PUT("/update/:id", controller.TrackUpdate)
		trk.POST("/avatar/:id", controller.TrackerAvatar)
		trk.GET("/geo/current/:id", controller.TrackerLastGeoPosition)
		trk.GET("/geo/history/:id", controller.TrackerHistory) //date_start date_end
	}

	zone := v1.Group("/geozones") // api для гео зон
	zone.Use(mFunc)
	{
		zone.GET("", controller.ZoneAllForUser)
		zone.POST("", controller.ZoneAdd)
		zone.DELETE("/:id", controller.ZoneDeleteById)
		zone.GET("/:id", controller.ZoneGetById)
		zone.PUT("/:id", controller.ZoneUpdate)
		zone.PUT("/:id/trackers", controller.ZoneSnapTrackList)
	}

	return r // initServer(":8070", r)
}

/*func createWatchApi() (*http.Server) {
	r := gin.New()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return initServer(":6666", r)
}*/

/*func dummyHandler(c *gin.Context) {
	c.AbortWithStatus(http.StatusGone)
}
*/
/*func initServer(port string, routes http.Handler) (*http.Server) {
	srv := &http.Server{
		Addr:    port,
		Handler: routes,
	}
	return srv
}
*/
