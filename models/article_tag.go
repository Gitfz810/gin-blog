package models

import "github.com/jinzhu/gorm"

type ArticleTag struct {
	ArticleID int
	TagID     int
}
// add tag to article
func addAritcleTag(articleID int, tagName string) error {
	tagID, err := GetTagIDByName(tagName)
	if err != nil {
		return err
	}

	articleTag := ArticleTag{
		ArticleID: articleID,
		TagID:     tagID,
	}

	if err = db.Create(&articleTag).Error; err != nil {
		return err
	}
	return nil
}

// delete tag from article
func deleteArticleTag(articleID int, tagName string) error {
	tagID, err := GetTagIDByName(tagName)
	if err != nil {
		return err
	}

	err = db.Where("article_id = ? AND tag_id = ?", articleID, tagID).Delete(&ArticleTag{}).Error
	if err != nil {
		return err
	}
	return nil
}

// ensure tags from article
func tags(articleID int) ([]string, error) {
	var articleTags []*ArticleTag
	err := db.Select("tag_id").Where("article_id = ?", articleID).Find(&articleTags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	res := make([]string, len(articleTags))
	for _, at := range articleTags {
		name, err := GetTagNameByID(at.TagID)
		if err != nil {
			return nil, err
		}
		res = append(res, name)
	}
	return res, nil
}
