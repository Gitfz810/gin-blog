package util

import "gin-blog/pkg/setting"

// Setup init the util
func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}
