package backend

import (
	"encoding/json"
	"errors"
	"github.com/Ryan0520/go-mmall/models"
	"github.com/Ryan0520/go-mmall/pkg/e"
	"github.com/Ryan0520/go-mmall/pkg/gredis"
	mmallCache "github.com/Ryan0520/go-mmall/service/cache"
	log "github.com/sirupsen/logrus"
)

type Product struct {
	P *models.Product

	OrderBy  string `json:"order_by"`
	Keyword  string `json:"keyword"`
	PageNum  int
	PageSize int
}

func (p *Product) Count() (int, error) {
	return models.GetProductTotalCount()
}

func (p *Product) SetSaleStatus() error {
	product := models.Product{
		Model:  models.Model{ID: p.P.ID},
		Status: p.P.Status,
	}
	return product.Update()
}

func (p *Product) Get() (*models.Product, error) {
	var cacheProduct *models.Product

	cache := mmallCache.Product{ID: p.P.ID}
	key := cache.GetProductKey()
	if gredis.Exist(key) {
		data, err := gredis.Get(key)
		if err != nil {
			log.Error(err)
		} else {
			json.Unmarshal(data, &cacheProduct)
			return cacheProduct, nil
		}
	}
	product, err := models.SelectProductById(p.P.ID)
	if err != nil {
		return nil, err
	}
	if product.ID == 0 {
		log.Error("产品 productId: ? 不存在", product.ID)
		return nil, errors.New(e.GetMsg(e.ERROR_NOT_EXIST_PRODUCT))
	}
	gredis.Set(key, product, 3600)
	return product, nil
}

func (p *Product) GetFilterOffSale() (*models.Product, error) {
	var cacheProduct *models.Product

	cache := mmallCache.Product{ID: p.P.ID}
	key := cache.GetProductKey()
	if gredis.Exist(key) {
		data, err := gredis.Get(key)
		if err != nil {
			log.Error(err)
		} else {
			json.Unmarshal(data, &cacheProduct)
			return cacheProduct, nil
		}
	}
	product, err := models.GetProductFilterOffSale(p.P.ID)
	if err != nil {
		return nil, err
	}
	if product.ID == 0 {
		log.Error("产品 productId: ? 不存在", product.ID)
		return nil, errors.New(e.GetMsg(e.ERROR_NOT_EXIST_PRODUCT))
	}
	gredis.Set(key, product, 3600)
	return product, nil
}

func (p *Product) GetAll() ([]*models.Product, error) {
	var (
		products, cacheProducts []*models.Product
	)

	id := 0
	if p.P != nil {
		id = p.P.ID
	}
	cache := mmallCache.Product{
		ID:       id,
		PageNum:  p.PageNum,
		PageSize: p.PageSize,
	}
	key := cache.GetProductsKey()
	if gredis.Exist(key) {
		data, err := gredis.Get(key)
		if err != nil {
			log.Error(err)
		} else {
			json.Unmarshal(data, &cacheProducts)
			return cacheProducts, nil
		}
	}
	products, err := models.GetProducts(p.PageNum, p.PageSize)
	if err != nil {
		return nil, err
	}
	gredis.Set(key, products, 3600)
	return products, nil
}

func (p *Product) GetAllFilterOffSale() ([]*models.Product, error) {
	var (
		products, _ []*models.Product
	)

	//id := 0
	//if p.P != nil {
	//	id = p.P.ID
	//}
	//cache := mmallCache.Product{
	//	ID:            id,
	//	Keyword:       p.Keyword,
	//	OrderBy:       p.OrderBy,
	//	PageNum:       p.PageNum,
	//	PageSize:      p.PageSize,
	//	FilterOffSale: true,
	//}
	//key := cache.GetProductsKey()
	//if gredis.Exist(key) {
	//	data, err := gredis.Get(key)
	//	if err != nil {
	//		log.Error(err)
	//	} else {
	//		json.Unmarshal(data, &cacheProducts)
	//		return cacheProducts, nil
	//	}
	//}
	products, err := models.GetProductsFilterOffSale(p.PageNum, p.PageSize, p.Keyword, p.OrderBy)
	if err != nil {
		return nil, err
	}
	//gredis.Set(key, products, 3600)
	return products, nil
}

func (p *Product) Search() ([]*models.Product, int, error) {
	return models.GetProductsByNameAndId(p.P.Name, p.P.ID, p.PageNum, p.PageSize)
}

func (p *Product) Save() error {
	return p.P.Save()
}
