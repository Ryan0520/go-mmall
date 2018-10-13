package models

import "github.com/jinzhu/gorm"

type Shipping struct {
	Model

	UserId           int    `json:"user_id"`
	ReceiverName     string `json:"receiver_name"`
	ReceiverPhone    string `json:"receiver_phone"`
	ReceiverMobile   string `json:"receiver_mobile"`
	ReceiverProvince string `json:"receiver_province"`
	ReceiverCity     string `json:"receiver_city"`
	ReceiverDistrict string `json:"receiver_district"`
	ReceiverAddress  string `json:"receiver_address"`
	ReceiverZip      string `json:"receiver_zip"`
}

func (s *Shipping)Save() (int, error) {
	if err := db.Create(s).Error; err != nil {
		return 0, err
	}

	return s.ID, nil
}

func (s *Shipping)Delete() error {
	return db.Delete(s).Error
}

func (s *Shipping)Update() error {
	data := make(map[string]interface{})
	if s.ID > 0 {
		data["id"] = s.ID
	}
	if s.UserId > 0 {
		data["user_id"] = s.UserId
	}
	if s.ReceiverName != "" {
		data["receiver_name"] = s.ReceiverName
	}
	if s.ReceiverPhone != "" {
		data["receiver_phone"] = s.ReceiverPhone
	}
	if s.ReceiverMobile != "" {
		data["receiver_mobile"] = s.ReceiverMobile
	}
	if s.ReceiverProvince != "" {
		data["receiver_province"] = s.ReceiverProvince
	}
	if s.ReceiverCity != "" {
		data["receiver_city"] = s.ReceiverCity
	}
	if s.ReceiverDistrict != "" {
		data["receiver_district"] = s.ReceiverDistrict
	}
	if s.ReceiverAddress != "" {
		data["receiver_address"] = s.ReceiverAddress
	}
	if s.ReceiverZip != "" {
		data["receiver_zip"] = s.ReceiverZip
	}
	return db.Model(s).Where("id = ?", s.ID).Updates(data).Error
}

func (s *Shipping)Get() error {
	return db.Model(&Shipping{}).First(s, "id = ? and user_id = ?", s.ID, s.UserId).Error
}

func GetShippingList(userId, pageNum, pageSize int) ([]*Shipping, int, error) {
	var list []*Shipping
	var totalCount int
	db.Model(&Shipping{}).Where("user_id = ?", userId).Count(&totalCount)
	err := db.Where("user_id = ?", userId).Offset(pageNum).Limit(pageSize).Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	return list, totalCount, nil
}