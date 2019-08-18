package jwt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"gin-blog/pkg/e"
	"gin-blog/pkg/util"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token, err := c.Cookie("token")
		fmt.Println(token)
		if token == "" || err != nil {
			fmt.Println("111")
			code = e.ERROR_AUTH_TOKEN
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg": e.GetMsg(code),
				"data": data,
			})
			// 中断执行，直接返回
			c.Abort()
			return
		}
		// 执行下一步  NopCloser 返回一个包裹起给定 Reader r 的 ReadCloser
		//c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		c.Next()
	}
}
