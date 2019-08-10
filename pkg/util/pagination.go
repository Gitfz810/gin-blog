package util

import (
	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"

	"gin-blog/pkg/setting"
)

func GetPage(c *gin.Context) int {
	result := 0
	/*
	1. 从gin的context上下文中获取page参数 string类型
	2. com.StrTo().Int() 就是对类型转换的一个封装 用的就是 strconv.ParseInt(str, 10, 0)
	 */
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * setting.PageSize
	}
	return result
}
