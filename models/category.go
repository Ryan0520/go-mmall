package models

import "github.com/jinzhu/gorm"

type Category struct {
	Model

	ParentId  int    `json:"parent_id"`
	Name      string `json:"name"`
	Status    bool   `json:"status"`
	SortOrder int    `json:"sort_order"`
}

func ExistCategoryByParentId(parentId int) (bool, error) {
	var count int
	err := db.Where("id = ?", parentId).Find(&Category{}).Count(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func GetCategoryById(id int) (Category, error) {
	var category Category
	err := db.Where("id = ?", id).Find(&Category{}).Scan(&category).Error
	return category, err
}

func GetCategoriesByParentId(parentId int) ([]*Category, error) {
	var categories []*Category
	err := db.Where("parent_id = ?", parentId).Find(&categories).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return categories, nil
}

func (category *Category) Save() error {
	if err := db.Create(category).Error; err != nil {
		return err
	}
	return nil
}

func (category *Category) Update() error {
	if err := db.Model(&category).Updates(category).Error; err != nil {
		return err
	}
	return nil
}
