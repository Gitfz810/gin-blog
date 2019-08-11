package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"gin-blog/models"
	"gin-blog/pkg/e"
	"gin-blog/pkg/logging"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
)

// 获取单个文章
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 0, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	var data interface{}
	if ! valid.HasErrors() {
		if ok, _ := models.ExistArticleByID(id); ok {
			code = e.SUCCESS
			data, _ = models.GetArticleById(id)
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": data,
	})
}

// 获取多篇文章
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})

	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	/*var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId
		valid.Min(tagId, 0, "tag_id").Message("标签ID必须大于0")
	}*/

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		code = e.SUCCESS

		data["lists"], _ = models.GetArticles(util.GetPage(c), setting.PageSize, maps)
		data["total"], _ = models.GetArticleTotal(maps)
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": data,
	})
}

// 新增文章
func AddArticle(c *gin.Context) {
	var jsonObj map[string]interface{}
	info, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(info, &jsonObj)
	if err != nil {
		logging.Fatal(err.Error())
	}

	title := jsonObj["title"].(string)
	desc := jsonObj["desc"].(string)
	content := jsonObj["content"].(string)
	createdBy := jsonObj["created_by"].(string)
	state := int(jsonObj["state"].(float64))

	valid := validation.Validation{}
	valid.Required(title, "title").Message("标题不能为空")
	valid.MaxSize(title, 50, "title").Message("标题字符最长50个字符")
	valid.MaxSize(desc, 100, "desc").Message("简述最长100个字符")
	//valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		data := make(map[string]interface{})
		data["title"] = title
		data["desc"] = desc
		data["content"] = content
		data["created_by"] = createdBy
		data["state"] = state

		models.AddArticle(data)
		code = e.SUCCESS
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]interface{}),
	})
}

// 修改文章
func EditArticle(c *gin.Context) {
	var jsonObj map[string]interface{}
	info, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(info, &jsonObj)
	if err != nil {
		logging.Fatal(err.Error())
	}

	id := com.StrTo(c.Param("id")).MustInt()
	title := jsonObj["title"].(string)
	tagNames := jsonObj["tag_names"].([]interface{})
	desc := jsonObj["desc"].(string)
	content := jsonObj["content"].(string)
	updatedBy := jsonObj["updated_by"].(string)

	valid := validation.Validation{}
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
	valid.Required(updatedBy, "updated_by").Message("修改人不能为空")
	valid.MaxSize(updatedBy, 100, "updated_by").Message("修改人最长为100字符")

	res := make([]string, len(tagNames))
	for i, name := range tagNames {
		res[i] = name.(string)
	}

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		if ok, _ := models.ExistArticleByID(id); ok {
			data := make(map[string]interface{})
			if title != "" {
				data["title"] = title
			}
			if desc != "" {
				data["desc"] = desc
			}
			if content != "" {
				data["content"] = content
			}
			data["updated_by"] = updatedBy
			models.UpdateTags(id, res)
			models.EditArticle(id, data)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]string),
	})
}

// 删除文章
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		if ok, _ := models.ExistArticleByID(id); ok {
			models.DeleteArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]string),
	})
}
