package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

const (
	OnSale = 1
	OffSale = 2
)

const (
	PRICE_DESC = "price desc"
	PRICE_ASC = "price asc"
)

type Product struct {
	Model

	CategoryId int `json:"category_id"`
	Name string `json:"name"`
	SubTitle string `json:"sub_title"`
	MainImage string `json:"main_image"`
	SubImages string `json:"sub_images"`
	Detail string `json:"detail"`
	Price int `json:"price"` // 单位是分
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

func SelectProductById(id int) (*Product, error)  {
	var product Product
	err := db.First(&product, "id = ?", id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &product, nil
}

func GetProductFilterOffSale(id int) (*Product, error)  {
	var product Product
	err := db.First(&product, "id = ? and status = ?", id, OnSale).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &product, nil
}

func GetProductsFilterOffSale(pageNum, pageSize int, keyword, orderBy string) ([]*Product, error) {
	var products []*Product
	var err error
	if len(keyword) == 0 {
		err = db.Where("status = ?", OnSale).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&products).Error
	} else {
		err = db.Offset(pageNum).Limit(pageSize).Order(orderBy).Where("name LIKE ? and status = ?", fmt.Sprintf("%%%s%%", keyword), OnSale).Find(&products).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return products, nil
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

func UpdateProductStock(id, stock int) error {
  	return db.Find(&Product{}, "id = ?", id).Updates(map[string]interface{}{
		"stock": stock,
	}).Error
}