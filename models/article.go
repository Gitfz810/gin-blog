package models

import "github.com/jinzhu/gorm"

type Article struct {
	Model

	TagID int `gorm:"index" json:"tag_id"`
	Tag   Tag `json:"tag"`  // 1 对 1

	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ? AND deleted_on is NULL", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if article.ID > 0 {
		return true, nil
	}
	return false, nil
}

func GetArticleTotal(maps interface{}) (count int, err error) {
	if err = db.Model(&Article{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return
}

func GetArticles(pageNum, pageSize int, maps interface{}) (article []*Article, err error) {
	// gorm默认不会查询外键对象，如果想把结构体字段的内容也查询出来，可以使用Preload函数预加载这个结构体
	err = db.Preload("Tag").Where(maps).Where("deleted_on is NULL").Offset(pageNum).Limit(pageSize).Find(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return
}

func GetArticleById(id int) (article *Article, err error) {
	/*// 更具id获取对应文章
	db.Where("id=?", id).First(&article)
	// 更具关联关系查找文章拥有的tag
	db.Model(&article).Related(&article.Tag)*/
	err = db.Preload("Tag").Where("id = ? AND deleted_on is NULL", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return
}

func EditArticle(id int, data interface{}) error {
	if err := db.Model(&Article{}).Where("id = ? AND deleted_on is NULL", id).Update(data).Error; err != nil {
		return err
	}

	return nil
}

func AddArticle(data map[string]interface{}) error {
	article := Article{
		TagID: data["tag_id"].(int),
		Title: data["title"].(string),
		Desc: data["desc"].(string),
		Content: data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State: data["state"].(int),
	}

	if err := db.Create(&article).Error; err != nil {
		return err
	}
	return nil
}

func DeleteArticle(id int) error {
	if err := db.Debug().Where("id = ?", id).Delete(Article{}).Error; err != nil {
		return err
	}
	return nil
}

/*func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", util.GetTime())

	return nil
}

func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", util.GetTime())

	return nil
}*/
