package portal

import (
	"errors"
	"github.com/Ryan0520/go-mmall/models"
	"github.com/Ryan0520/go-mmall/pkg/e"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"strings"
)

const (
	LIMIT_NUM_SUCCESS = "LIMIT_NUM_SUCCESS"
	LIMIT_NUM_FAIL    = "LIMIT_NUM_FAIL"
)

type CartProductVo struct {
	ID                int    `json:"id"`
	UserId            int    `json:"user_id"`
	ProductId         int    `json:"product_id"`
	Quantity          int    `json:"quantity"`
	ProductName       string `json:"product_name"`
	ProductMainImage  string `json:"product_main_image"`
	ProductPrice      int    `json:"product_price"` // 单位是分
	ProductStatus     int    `json:"product_status"`
	ProductTotalPrice int    `json:"product_total_price"` // 单位是分
	ProductStock      int    `json:"product_stock"`
	ProductChecked    int    `json:"product_checked"`
	LimitQuantity     string `json:"limit_quantity"`
}

type CartVo struct {
	Carts      []*CartProductVo `json:"carts"`
	TotalPrice int              `json:"total_price"`
	AllChecked bool             `json:"all_checked"`
}

func (c *CartProductVo) GetCartList() (*CartVo, error) {
	return c.getCartVO()
}

func (c *CartProductVo) AddCartProduct(productId, count int) (*CartVo, error) {
	cart, err := models.SelectCartWithUserIdAndProductId(c.UserId, productId)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("SelectCartWithUserIdAndProductId Error: %v", err)
		return nil, err
	}

	// 没有该商品的购物车记录，需要添加
	if cart == nil || err == gorm.ErrRecordNotFound {
		p, err := models.GetProduct(productId)
		if err != nil{
			log.Errorf("产品 productId: %d 不存在", productId)
			return nil, err
		}
		if p.ID == 0 {
			log.Errorf("产品 productId: %d 不存在", productId)
			return nil, errors.New(e.GetMsg(e.ERROR_NOT_EXIST_PRODUCT))
		}
		cart = &models.Cart{
			UserId: c.UserId,
			ProductId: productId,
			Quantity: count,
			Checked: models.Checked,
		}
		err = cart.Create()
		if err != nil {
			log.Error("create cart err: ", err)
			return nil, err
		}
	} else {
		count = cart.Quantity + count
		models.UpdateCart(cart.ID, map[string]interface{}{
			"quantity": count,
		})
	}
	return c.getCartVO()
}

func (c *CartProductVo)UpdateCartProductCount(productId, count int) (*CartVo, error) {
	cart, err := models.SelectCartWithUserIdAndProductId(c.UserId, productId)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("SelectCartWithUserIdAndProductId Error: %v", err)
		return nil, err
	}

	// 没有该商品的购物车记录，需要添加
	if cart == nil || err == gorm.ErrRecordNotFound {
		p, err := models.GetProduct(productId)
		if err != nil{
			log.Errorf("产品 productId: %d 不存在", productId)
			return nil, err
		}
		if p.ID == 0 {
			log.Errorf("产品 productId: %d 不存在", productId)
			return nil, errors.New(e.GetMsg(e.ERROR_NOT_EXIST_PRODUCT))
		}
		cart = &models.Cart{
			UserId: c.UserId,
			ProductId: productId,
			Quantity: count,
			Checked: models.Checked,
		}
		err = cart.Create()
		if err != nil {
			log.Error("create cart err: ", err)
			return nil, err
		}
	} else {
		models.UpdateCart(cart.ID, map[string]interface{}{
			"quantity": count,
		})
	}
	return c.getCartVO()
}

func (c *CartProductVo)DeleteCartProducts(productIds string) (*CartVo, error) {
	err := models.DeleteCartProductsByUserIdAndProductIds(c.UserId, strings.Split(productIds, ","))
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return c.getCartVO()
}

func (c *CartProductVo)SelectCartProduct(productId int) (*CartVo, error) {
	err := models.UpdateCartProductCheckedStatusByProductId(c.UserId, productId, models.Checked)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return c.getCartVO()
}

func (c *CartProductVo)UnSelectCartProduct(productId int) (*CartVo, error) {
	err := models.UpdateCartProductCheckedStatusByProductId(c.UserId, productId, models.UnChecked)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return c.getCartVO()
}

func (c *CartProductVo)GetCartProductCount() (int, error) {
	count, err := models.SelectCartProductCount(c.UserId)
	if err !=nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return count, err
}

func (c *CartProductVo)SelectAllCart() (*CartVo, error) {
	err := models.UpdateAllCartCheckedStatus(c.UserId, models.Checked)
	if err !=nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return c.getCartVO()
}

func (c *CartProductVo)UnSelectAllCart() (*CartVo, error) {
	err := models.UpdateAllCartCheckedStatus(c.UserId, models.UnChecked)
	if err !=nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return c.getCartVO()
}

func (c *CartProductVo)getCartVO() (*CartVo, error) {
	mCarts, err := models.GetCartsByUserID(c.UserId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var carts []*CartProductVo
	totalPrice := 0
	for index := 0; index < len(mCarts); index++ {
		mCart := mCarts[index]
		var cartProduct CartProductVo
		cartProduct.ID = mCart.ID
		cartProduct.UserId = mCart.UserId
		cartProduct.ProductId = mCart.ProductId

		mProduct, err := models.GetProduct(mCart.ProductId)
		if err != nil {
			log.Error("获取product失败, productId：%d", mCart.ProductId)
		} else {
			if mProduct.ID > 0 {
				cartProduct.ProductMainImage = mProduct.MainImage
				cartProduct.ProductPrice = mProduct.Price
				cartProduct.ProductStock = mProduct.Stock
				cartProduct.ProductName = mProduct.Name
				cartProduct.ProductStatus = mProduct.Status
				buyLimitCount := 0
				if mProduct.Stock >= mCart.Quantity {
					buyLimitCount = mCart.Quantity
					cartProduct.LimitQuantity = LIMIT_NUM_SUCCESS
				} else {
					buyLimitCount = mProduct.Stock
					cartProduct.LimitQuantity = LIMIT_NUM_FAIL
					models.UpdateCart(mCart.ID, map[string]interface{}{
						"quantity": buyLimitCount,
					})
				}
				cartProduct.Quantity = buyLimitCount
				cartProduct.ProductTotalPrice = mCart.Quantity * mProduct.Price
				cartProduct.ProductChecked = mCart.Checked
				if mCart.Checked == models.Checked {
					totalPrice += cartProduct.ProductTotalPrice
				}
				carts = append(carts, &cartProduct)
			} else {
				log.Error("productID == 0")
			}
		}
	}
	unSelectCount, err := models.SelectCartProductCheckedStatusByUserId(c.UserId)
	if err != nil {
		log.Errorf("SelectCartProductCheckedStatusByUserId err: %v", err)
	}

	cartVo := &CartVo{
		Carts:      carts,
		TotalPrice: totalPrice,
		AllChecked: unSelectCount == 0,
	}
	return cartVo, nil
}
