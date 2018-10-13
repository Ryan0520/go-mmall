package portal

import (
	"github.com/Ryan0520/go-mmall/models"
)

type Shipping struct {
	ID               int    `json:"id"`
	UserId           int    `json:"user_id"`
	ReceiverName     string `json:"receiver_name"`
	ReceiverPhone    string `json:"receiver_phone"`
	ReceiverMobile   string `json:"receiver_mobile"`
	ReceiverProvince string `json:"receiver_province"`
	ReceiverCity     string `json:"receiver_city"`
	ReceiverDistrict string `json:"receiver_district"`
	ReceiverAddress  string `json:"receiver_address"`
	ReceiverZip      string `json:"receiver_zip"`

	PageNum  int
	PageSize int
}

func (s *Shipping)Add() (int, error) {
	modelS := models.Shipping{
		Model: models.Model{
			ID:s.ID,
		},
		UserId: s.UserId,
		ReceiverName: s.ReceiverName,
		ReceiverPhone: s.ReceiverPhone,
		ReceiverMobile: s.ReceiverMobile,
		ReceiverProvince: s.ReceiverProvince,
		ReceiverCity: s.ReceiverCity,
		ReceiverAddress: s.ReceiverAddress,
		ReceiverDistrict: s.ReceiverDistrict,
		ReceiverZip: s.ReceiverZip,
	}
	return modelS.Save()
}

func (s *Shipping)Delete() error {
	modelS := models.Shipping{
		Model: models.Model{
			ID: s.ID,
		},
	}
	return modelS.Delete()
}

func (s *Shipping)Update() error {
	modelS := models.Shipping{
		Model: models.Model{
			ID:s.ID,
		},
		UserId: s.UserId,
		ReceiverName: s.ReceiverName,
		ReceiverPhone: s.ReceiverPhone,
		ReceiverMobile: s.ReceiverMobile,
		ReceiverProvince: s.ReceiverProvince,
		ReceiverCity: s.ReceiverCity,
		ReceiverAddress: s.ReceiverAddress,
		ReceiverDistrict: s.ReceiverDistrict,
		ReceiverZip: s.ReceiverZip,
	}
	return modelS.Update()
}

func (s *Shipping)Get() (*models.Shipping, error) {
	modelS := models.Shipping{
		Model: models.Model{
			ID: s.ID,
		},
		UserId:s.UserId,
	}
	err := modelS.Get()
	if err != nil {
		return nil, err
	}
	return &modelS, nil
}

func (s *Shipping)GetList() ([]*models.Shipping, int, error) {
	return models.GetShippingList(s.UserId, s.PageNum, s.PageSize)
}