package router

import (
	"crypto/tls"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/docs"
	"github.com/ppoonk/AirGo/global"
	middleware "github.com/ppoonk/AirGo/router/middleware"
	"github.com/ppoonk/AirGo/web"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	"os"
	"strconv"
	"sync"
)

type GinRouter struct {
	Router *gin.Engine
}

func NewGinRouter() *GinRouter {
	return &GinRouter{Router: nil}
}

func (g *GinRouter) InitRouter() {
	gin.SetMode(gin.ReleaseMode)
	var writer io.Writer
	if global.Config.SystemParams.Mode == "dev" {
		writer = os.Stdout
	} else {
		writer = io.Discard //关闭控制台输出
	}
	gin.DefaultWriter = writer
	g.Router = gin.Default()

	// 静态资源路由，以 / 开头，且不是 /api 开头的，都被当成静态资源，静态资源目录为 项目/server/web/web

	// targetPtah=web 是embed和web文件夹的相对路径
	g.Router.Use(middleware.Serve("/", middleware.EmbedFolder(web.Static, "web")))
	g.Router.Use(middleware.Cors(), middleware.Recovery())

	// 业务api路由组
	apiRouter := g.Router.Group("/api")
	apiRouter.Use(middleware.DomainAndAPI())

	// swagger 路由 例如：http://localhost:8899/api/swagger/index.html
	docs.SwaggerInfo.BasePath = ""
	swaggerRouter := apiRouter.Group("/swagger")
	{
		swaggerRouter.GET("/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	// 注册业务路由
	g.InitAdminRouter(apiRouter)
	g.InitUserRouter(apiRouter)
	g.InitPublicRouter(apiRouter)
}

func (g *GinRouter) Start() {
	w := sync.WaitGroup{}
	w.Add(2)
	go func() {
		err := endless.ListenAndServe(":"+strconv.Itoa(global.Config.SystemParams.HTTPPort), g.Router)
		if err != nil {
			global.Logrus.Error("listen: %s", err)
		}
		w.Done()
	}()
	go func() {
		_, err := tls.LoadX509KeyPair("./air.cer", "./air.key") //先验证证书，否则endless fork进程时会空指针panic
		if err == nil {
			err = endless.ListenAndServeTLS(":"+strconv.Itoa(global.Config.SystemParams.HTTPSPort), "./air.cer", "./air.key", g.Router)
			if err != nil {
				global.Logrus.Error("listen: %s", err)
			}
		}
		w.Done()
	}()
	w.Wait()

	//syscall.SIGHUP 将触发重启; syscall.SIGINT, syscall.SIGTERM 并将触发服务器关闭（它将完成运行请求)。https://github.com/fvbock/endless
	// windows下使用endless报错：undefined: syscall.SIGUSR1
	// 解决办法：https://github.com/golang/go/pull/51541/commits/1a70eff904944dde32213a27d3357c7380067fd0

	global.Logrus.Info("Server stop")
	os.Exit(0)
}
