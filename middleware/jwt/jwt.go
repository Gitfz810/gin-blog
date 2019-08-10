package jwt

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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
		tmp := make(map[string]interface{})
		body, _ := ioutil.ReadAll(c.Request.Body)
		json.Unmarshal(body, &tmp)
		token, ok := tmp["token"]
		if token == "" || !ok {
			code = e.INVALID_PARAMS
		} else {
			claims, err := util.ParseToken(token.(string))
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
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		c.Next()
	}
}
