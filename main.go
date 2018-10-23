package main

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gobeacon/controller"
	_ "gobeacon/docs"
	"gobeacon/watch"
	"log"
	"net/http"
	"sync"
	"time"
)

// @title Swagger API
// @version 1.0
// @description This is api for beacon services
// @contact.name API Support
// @contact.email paradoxfm@mail.ru
func main() {
	//db.InitConnect()
	//gin.SetMode(gin.ReleaseMode)
	startTcpServer()
	// если пойдет большая нагрузка, то надо распилить на отдельные приложения
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

func startTcpServer() {
	srv := watch.Server{
		Addr:         ":6666",
		IdleTimeout:  150 * time.Second,
		MaxReadBytes: 10000,
	}
	go srv.ListenAndServe()
}

func createPhoneAdminApi() (*gin.Engine) {
	auth := controller.CreateAdminJWTMiddleware()
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	// 1 << 20  1 MiB -> ‭1_048_576‬, 8 << 20  8 MiB -> ‭8_388_608‬
	r.MaxMultipartMemory = 1 << 19 //0.5 MiB
	v1 := r.Group("/api/v1")       // api первой версии
	mFunc := auth.MiddlewareFunc()

	// новое api для смартфонов на том же порту
	bgr := v1.Group("", gin.BasicAuthForRealm(gin.Accounts{
		"heart349023": "s156EzI07820CtsfJhu",
	}, "phone connector"))
	bgr.POST("/heartbeat", controller.HeartbeatPhone)

	tst := v1.Group("/test")
	tst.Use(mFunc)
	tst.GET("/push", controller.TestPush)
	//tst.GET("/updtrack", controller.TestTrack)
	sys := v1.Group("")
	sys.Use(mFunc)
	{
		sys.GET("/avatar/:id", controller.GetAvatar)
	}
	usr := v1.Group("/users") // api для пользователей
	{
		usr.POST("/signup", controller.UserCreate)
		usr.POST("/login", auth.LoginHandler)
		usr.POST("/password/reset", controller.UserResetPassword)
		me := usr.Group("/me")
		me.Use(mFunc)
		{
			me.GET("", controller.UserGetProfile)
			me.PUT("/password", controller.UserChangePassword)
			me.PUT("/push", controller.UserUpdatePushId)
			me.PUT("/avatar", controller.UserUpdateAvatar)
			//me.PUT("/refresh", auth.RefreshHandler)
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
		trk.PUT("/avatar/:id", controller.TrackerAvatar)
		//trk.GET("/geo/current/:id", controller.TrackerLastGeoPosition)
		trk.POST("/geo/history", controller.TrackerHistory)
	}

	zone := v1.Group("/zone") // api для гео зон
	zone.Use(mFunc)
	{
		zone.GET("/all", controller.ZoneAllForUser)
		zone.POST("/save", controller.ZoneCreate)
		zone.DELETE("/delete/:id", controller.ZoneDeleteById)
		zone.GET("/find/:id", controller.ZoneGetById)
		zone.PUT("/update/:id", controller.ZoneUpdate)
		zone.PUT("/snap/:id", controller.ZoneSnapTrackList)
	}

	return r
}

//depreceated from 23.10.2018
func createPhoneApi() (*gin.Engine) {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	v1 := r.Group("/api/v1", gin.BasicAuthForRealm(gin.Accounts{
		"heart349023": "s156EzI07820CtsfJhu",
	}, "phone connector"))
	v1.POST("/heartbeat", controller.HeartbeatPhone)
	return r
}

func createSwaggerApi() (*gin.Engine) {
	r := gin.New()
	r.Use(gin.Recovery())
	authorized := r.Group("", gin.BasicAuth(gin.Accounts{
		"admin": "password",
	}))
	// документация по сервисам /swagger/index.html
	authorized.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
