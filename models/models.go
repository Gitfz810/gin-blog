package models

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
)
// 全局变量 在整个models包内，所有文件均可以直接使用db
var db *gorm.DB

type Model struct {
	ID        int        `gorm:"primary_key" json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	//gorm.Model
}

func init() {
	var (
		err error
		connTimeout int
		dbType, dbName, user, password, host, tablePrefix string
	)

	dbInfo, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	dbType = dbInfo.Key("TYPE").String()
	dbName = dbInfo.Key("NAME").String()
	user = dbInfo.Key("USER").String()
	password = dbInfo.Key("PASSWORD").String()
	host = dbInfo.Key("HOST").String()
	tablePrefix = dbInfo.Key("TABLE_PREFIX").String()
	connTimeout, _ = dbInfo.Key("CONN_TIMEOUT").Int()

	// 正确的处理 time.Time ，需要包含 parseTime 参数 loc 指定时区
	db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName))
	if err != nil {
		log.Println(err)
	}

	// 更改默认表名 对表名增加前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}
	// 禁用表名复数形式
	db.SingularTable(true)
	// 设置空闲连接
	db.DB().SetMaxIdleConns(5)
	// 设置最大连接数
	db.DB().SetMaxOpenConns(30)
	// 设置连接过期时间
	db.DB().SetConnMaxLifetime(time.Duration(connTimeout) * time.Hour)
	// 替换默认回调函数
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	// 显示执行的SQL语句
	db.LogMode(true)
}

func CloseDB() {
	defer db.Close()
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if ! scope.HasError() {
		nowTime := util.GetTime()
		if createTimeField, ok := scope.FieldByName("CreatedAt"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("UpdatedAt"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	// 查询是否有指定`update_column`字段的column
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("UpdatedAt", util.GetTime())
	}
}

// deleteCallback will set `deletedOn` when delete data
func deleteCallback(scope *gorm.Scope) {
	if ! scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedAt")
		if ! scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(util.GetTime()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
				)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
