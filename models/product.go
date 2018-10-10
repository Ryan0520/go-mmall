package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

const (
	OnSale = 1
	OffSale = 2
)

type Product struct {
	Model

	CategoryId int `json:"category_id"`
	Name string `json:"name"`
	SubTitle string `json:"sub_title"`
	MainImage string `json:"main_image"`
	SubImages string `json:"sub_images"`
	Detail string `json:"detail"`
	Price float64 `json:"price"`
	Stock int `json:"stock"`
	Status int `json:"status"`
}

func GetProductTotalCount() (int, error) {
	var count int
	if err := db.Model(&Product{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetProduct(id int) (*Product, error)  {
	var product Product
	err := db.First(&product, "id = ?", id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &product, nil
}

func GetProducts(pageNum, pageSize int) ([]*Product, error) {
	var products []*Product
	err := db.Offset(pageNum).Limit(pageSize).Find(&products).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return products, nil
}

func GetProductsByNameAndId(productName string, productId int, pageNum, pageSize int) ([]*Product, int, error) {
	var products []*Product
	totalCount := 0
	if len(productName) > 0 && productId > 0 {
		db.Where("name LIKE ? AND id = ?", fmt.Sprintf("%%%s%%", productName), productId).Find(&products).Count(&totalCount)
		err := db.Where("name LIKE ? AND id = ?", fmt.Sprintf("%%%s%%", productName), productId).Offset(pageNum).Limit(pageSize).Find(&products).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, totalCount, err
		}
		return products, totalCount, nil
	}

	if len(productName) > 0 {
		db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", productName)).Find(&products).Count(&totalCount)
		err := db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", productName)).Offset(pageNum).Limit(pageSize).Find(&products).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, totalCount, err
		}
		return products, totalCount, nil
	}

	if productId > 0 {
		db.Where("id = ?", productId).Find(&products).Count(&totalCount)
		err := db.Where("id = ?", productId).Offset(pageNum).Limit(pageSize).Find(&products).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, totalCount, err
		}
		return products, totalCount, nil
	}

	return products, totalCount, nil
}

func (p *Product) Save() error {
	if err := db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Product) Update() error {
	if err := db.Model(p).Updates(p).Error; err != nil {
		return err
	}
	return nil
}