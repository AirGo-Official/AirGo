package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/docs"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/middleware"
	"github.com/ppoonk/AirGo/web"
	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var Router *gin.Engine

func InitRouter() {
	gin.SetMode(gin.ReleaseMode) //ReleaseMode TestMode DebugMode
	Router = gin.Default()
	// targetPtah=web 是embed和web文件夹的相对路径
	Router.Use(middleware.Serve("/", middleware.EmbedFolder(web.Static, "web")))
	Router.Use(middleware.Cors(), middleware.Recovery())

	//api路由
	RouterGroup := Router.Group("/api")

	//swagger 路由
	docs.SwaggerInfo.BasePath = ""
	swaggerRouter := RouterGroup.Group("/swagger").Use(middleware.ParseJwt(), middleware.Casbin())
	{
		swaggerRouter.GET("/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	//注册路由
	InitAdminRouter(RouterGroup)
	InitUserRouter(RouterGroup)
	InitPublicRouter(RouterGroup)

}

func ListenAndServe() {

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(global.Config.SystemParams.HTTPPort),
		Handler: Router,
	}
	srvTls := &http.Server{
		Addr:    ":" + strconv.Itoa(global.Config.SystemParams.HTTPSPort),
		Handler: Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logrus.Fatalf("listen: %s\n", err)
		}
	}()
	go func() {
		if err := srvTls.ListenAndServeTLS("./air.cer", "./air.key"); err != nil && err != http.ErrServerClosed {
			global.Logrus.Error("tls listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	global.Logrus.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.Logrus.Fatalf("Server Shutdown:", err)
	}
	if err := srvTls.Shutdown(ctx); err != nil {
		global.Logrus.Fatalf("Server Shutdown:", err)
	}
	global.Logrus.Info("Server exiting")
}
