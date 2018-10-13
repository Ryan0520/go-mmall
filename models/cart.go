package models

import "github.com/jinzhu/gorm"

const (
	Checked = 1
	UnChecked = 0
)

type Cart struct {
	Model

	UserId    int `json:"user_id"`
	ProductId int `json:"product_id"`
	Quantity  int `json:"quantity"`
	Checked   int `json:"checked"`
}

func GetCartsByUserID(userId int) ([]*Cart, error)  {
	var carts []*Cart
	err := db.Find(&carts, "user_id = ?", userId).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return carts, nil
}

func UpdateCart(id int, data map[string]interface{}) error {
	if err := db.Model(&Cart{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func SelectCartWithUserIdAndProductId(userId, productId int) (*Cart, error) {
	var cart Cart
	if err := db.Find(&cart, "user_id = ? and product_id = ?", userId, productId).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func SelectCartProductCheckedStatusByUserId(userId int) (int, error) {
	var count int
	if err := db.Find(&Cart{}, "user_id = ? and checked = 0", userId).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func DeleteCartProductsByUserIdAndProductIds(userId int, productIds []string) error {
	return db.Delete(Cart{}, "user_id = ? and product_id in (?)", userId, productIds).Error
}

func UpdateCartProductCheckedStatusByProductId(userId, productId, checked int) error {
	data := make(map[string]interface{})
	data["checked"] = checked
	return db.Model(&Cart{}).Where("user_id = ? and product_id = ?", userId, productId).Updates(data).Error
}

func UpdateAllCartCheckedStatus(userId, checked int) error {
	data := make(map[string]interface{})
	data["checked"] = checked
	return db.Model(&Cart{}).Where("user_id = ?", userId).Updates(data).Error
}

func SelectCartProductCount(userId int) (int, error)  {
	var count int
	err := db.Raw("select IFNULL(sum(quantity),0) as count from mmall_cart where user_id = ? and deleted_at is null", userId).Row().Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *Cart) Save() error  {
	if err := db.Create(c).Error; err != nil {
		return err
	}
	return nil
}