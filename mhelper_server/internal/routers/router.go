package routers

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"mhelper_server/global"
	"mhelper_server/internal/middleware"
	"mhelper_server/internal/routers/api"
	v1 "mhelper_server/internal/routers/api/v1"
	"mhelper_server/pkg/limiter"
	"net/http"
	"time"
)

// 设置限流规则
var methodLimiters = limiter.NewMethodLimiter().AddBuckets(
	limiter.BucketRule{
		Key:          "/auth",
		FillInterval: time.Second,
		Capacity:     10, // 容量 TODO ？？？
		Quantum:      10, // 大小 ？？？
	},
)

func NewRouter() *gin.Engine {
	r := gin.New()
	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery()) // gin 默认的两个中间件
	} else {
		r.Use(middleware.AccessLog()) // 访问日志
		r.Use(middleware.Recovery())  // Recovery
	}

	r.Use(middleware.Tracing())                                               // 链路追踪
	r.Use(middleware.RateLimiter(methodLimiters))                             // 限流器
	r.Use(middleware.ContextTimeout(global.AppSetting.DefaultContextTimeout)) // 超时处理
	r.Use(middleware.Translations())                                          // 多语言支持

	article := v1.NewArticle()
	tag := v1.NewTag()
	upload := api.NewUpload()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.GET("/debug/vars", api.Expvar)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload/file", upload.UploadFile)
	r.POST("/auth", api.GetAuth)

	// 实现让外部请求本项目 HTTP Server 时同时提供静态资源的访问，需要设置文件服务去提供静态资源的访问
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))
	apiv1 := r.Group("/api/v1")
	apiv1.Use() //middleware.JWT()
	{
		// 创建标签
		apiv1.POST("/tags", tag.Create)
		// 删除指定标签
		apiv1.DELETE("/tags/:id", tag.Delete)
		// 更新指定标签
		apiv1.PUT("/tags/:id", tag.Update)
		// 获取标签列表
		apiv1.GET("/tags", tag.List)

		// 创建文章
		apiv1.POST("/articles", article.Create)
		// 删除指定文章
		apiv1.DELETE("/articles/:id", article.Delete)
		// 更新指定文章
		apiv1.PUT("/articles/:id", article.Update)
		// 获取指定文章
		apiv1.GET("/articles/:id", article.Get)
		// 获取文章列表
		apiv1.GET("/articles", article.List)
	}

	return r
}
