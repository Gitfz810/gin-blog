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

// 获取多个文章标签
func GetTags(c *gin.Context) {
	// c.Query() 获取url 传参 ?name=test&state=1  c.DefaultQuery() 支持设置一个默认值
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		// state, _ = arg.(Int)
		state = com.StrTo(arg).MustInt()
		maps["state" ] = state
	}

	code := e.SUCCESS

	data["lists"], _ = models.GetTags(util.GetPage(c), setting.PageSize, maps)

	data["total"], _ = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": data,
	})
}

// 新增文章标签
func AddTags(c *gin.Context) {
	var jsonObj map[string]interface{}
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &jsonObj)
	if err != nil {
		logging.Fatal(err.Error())
	}

	name := jsonObj["name"].(string)
	state := int(jsonObj["state"].(float64))
	createdBy := jsonObj["created_by"].(string)
	/*name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")*/

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	// 判断name state 和 createdBy 字段是否正确
	if ! valid.HasErrors() {
		// 判断 tag 表中是否已存在以 name 为名字的 tag
		if ok, _ := models.ExistTagByName(name); !ok {  // 不存在
			code = e.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			code = e.ERROR_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// 修改文章标签
func EditTags(c *gin.Context) {
	// c.Param 获取的是URL路径参数 c.Query 获取的是URL ?拼接后的参数
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许为0或1")
	}

	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		code = e.SUCCESS
		if ok, _ := models.ExistTagByID(id); ok {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			data["name"] = name
			if state != -1 {
				data["state"] = state
			}
			models.EditTag(id, data)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// 删除文章标签
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Query("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 0, "id").Message("id必须大于0")

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		code = e.SUCCESS
		if ok, _ := models.ExistTagByID(id); ok {
			models.DeleteTag(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": make(map[string]string),
	})
}