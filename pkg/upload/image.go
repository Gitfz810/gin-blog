package upload

import (
	"fmt"
	"gin-blog/pkg/file"
	"gin-blog/pkg/logging"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

// get image name 获取图片名称
func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)
	return fileName + ext
}

// get save path 获取图片路径
func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

// get full save path 获取图片完整路径
func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

// get the full access path 获取图片完整访问URL
func GetImageFullUrl(name string) string {
	return strings.Join([]string{setting.AppSetting.PrefixUrl, "/", GetImagePath(), name}, "")
}

// check image file ext 检查图片后缀
func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToLower(allowExt) == strings.ToLower(ext) {
			return true
		}
	}
	return false
}

// check image size 检查图片大小
func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}

	return size <= setting.AppSetting.ImageMaxSize
}

// check image file is exists or not 检查图片是否存在
func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = file.IsNotExistMkDir(strings.Join([]string{dir, "/", src}, ""))
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := file.CheckPermission(src)
	if perm {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}
