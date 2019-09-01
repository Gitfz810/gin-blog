package models

import (
	"github.com/jinzhu/gorm"

	"gin-blog/pkg/set"
)

type Article struct {
	Model

	//TagID     int    `gorm:"index" json:"tag_id"`
	Tag           []Tag  `gorm:"many2many:article_tag;" json:"tag"`  // 多对多

	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	CreatedBy     string `json:"created_by"`
	UpdatedBy     string `json:"updated_by"`
	State         int    `json:"state"`
}

func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ?", id).First(&article).Error
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
	err = db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return
}

func GetArticleByID(id int) (articles *Article, err error) {
	err = db.Preload("Tag").Where("id = ?", id).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return
}

func EditArticle(id int, data interface{}) error {
	if err := db.Model(&Article{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func AddArticle(data map[string]interface{}) error {
	article := Article{
		Title: data["title"].(string),
		Desc: data["desc"].(string),
		Content: data["content"].(string),
		CoverImageUrl: data["cover_image_url"].(string),
		CreatedBy: data["created_by"].(string),
		State: data["state"].(int),
	}

	if err := db.Create(&article).Error; err != nil {
		return err
	}
	return nil
}

func UpdateTags(articleID int, 	tagNames []string) error {
	// 查询已存在的 tag
	existTags, err := tags(articleID)
	if err != nil {
		return err
	}

	// 确保需要插入 tag 存在
	for _, name := range tagNames {
		if ok, _ := ExistTagByName(name); !ok {
			err := AddTag(name, 1, "")
			if err != nil {
				return err
			}
		}
	}

	updateTagNames := set.New(tagNames)
	existsTagNames := set.New(existTags)

	needAddNames := updateTagNames.Minus(existsTagNames)
	for _, name := range needAddNames.SortList() {
		err := addAritcleTag(articleID, name)
		if err != nil {
			return err
		}
	}
	needDeleteNames := existsTagNames.Minus(updateTagNames)
	for _, name := range needDeleteNames.SortList() {
		err := deleteArticleTag(articleID, name)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteArticle(id int) error {
	if err := db.Where("id = ?", id).Delete(Article{}).Error; err != nil {
		return err
	}
	return nil
}
// 硬删除 使用Unscoped() GORM的约定
func CleanAllArticle() error {
	if err := db.Unscoped().Delete(&Article{}).Error; err != nil {
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
