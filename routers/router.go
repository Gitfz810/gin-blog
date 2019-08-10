package routers

import (
	"github.com/gin-gonic/gin"

	"gin-blog/middleware/jwt"
	"gin-blog/pkg/setting"
	"gin-blog/routers/api"
	"gin-blog/routers/api/v1"
)

func InitRouter() *gin.Engine {
	// 关闭控制台输出
	//gin.DisableConsoleColor()
	// 将输出信息写入日志
	//gin.DefaultWriter = io.MultiWriter(logging.F)
	r := gin.New()
	// 加载gin自带的logger和recovery中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// 通过setting.RunMode设置运行模式
	gin.SetMode(setting.RunMode)

	r.GET("/auth/", api.GetAuth)
	// test
	/*r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})*/

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		// 获取标签列表 path /api/v1/tags
		apiv1.GET("/tags/", v1.GetTags)
		// 新建标签
		apiv1.POST("/tags/", v1.AddTags)
		// 更新执行标签
		apiv1.PUT("/tags/:id/", v1.EditTags)
		// 删除标签
		apiv1.DELETE("/tags/:id/", v1.DeleteTag)

		// 获取文章列表
		apiv1.GET("/articles/", v1.GetArticles)
		// 获取指定文章
		apiv1.GET("/articles/:id/", v1.GetArticle)
		// 新建文章
		apiv1.POST("/articles/", v1.AddArticle)
		// 更新指定文章
		apiv1.PUT("/articles/:id/", v1.EditArticle)
		// 删除指定文章
		apiv1.DELETE("/articles/:id/", v1.DeleteArticle)
	}

	return r
}
