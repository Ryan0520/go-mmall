package portal

import (
	"fmt"
	"github.com/Ryan0520/go-mmall/models"
	"github.com/Ryan0520/go-mmall/pkg/e"
	"github.com/Ryan0520/go-mmall/pkg/setting"
	ali "github.com/Ryan0520/go-mmall/service/alipay"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/core/errors"
	log "github.com/sirupsen/logrus"
	"github.com/smartwalle/alipay"
	"math/rand"
	"strconv"
	"time"
)

type Order struct {
	OrderNo    int64
	UserId     int
	ShippingId int
	PageNum    int
	PageSize   int
}

func (o *Order) PayOrder() (string, error) {
	modelO, err := models.SelectOrderByUserIdAndOrderNo(o.UserId, o.OrderNo)
	if err != nil {
		return "", err
	}
	if modelO.ID <= 0 {
		return "", errors.New(e.GetMsg(e.ERROR_NOT_EXIST_ORDER))
	}
	if modelO.Status >= models.ORDER_PAID {
		log.Infof("订单号:%v 已付款", modelO.OrderNo)
		return "", errors.New("订单已支付")
	}

	if modelO.Status == models.ORDER_PAY_CANCEL {
		log.Infof("订单号:%v 已取消", modelO.OrderNo)
		return "", errors.New("订单已取消")
	}

	log.Infoln("========== TradePagePay ==========")
	var p = alipay.AliPayTradePagePay{}
	p.NotifyURL = setting.AlipaySetting.NotifyUrl
	p.ReturnURL = setting.AlipaySetting.ReturnUrl
	p.Subject = "订单" + fmt.Sprintf("%v", o.OrderNo) + "购买商品共" + strconv.Itoa(modelO.Payment/100) + "元"
	p.OutTradeNo = fmt.Sprintf("%v", o.OrderNo)
	p.TotalAmount = strconv.Itoa(modelO.Payment / 100)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	url, err := ali.Client.TradePagePay(p)
	if err != nil {
		log.Error(err)
		return "", err
	}
	log.Infoln(url)
	if len(url.String()) > 0 {
		log.Infof("订单号:%v 支付宝支付预创建成功 支付url:%s", o.OrderNo, url.String())
	}
	return url.String(), nil
}

func (o *Order) QueryOrderPayStatus() (bool, error) {
	modelO, err := models.SelectOrderByUserIdAndOrderNo(o.UserId, o.OrderNo)
	if err != nil {
		return false, err
	}
	if modelO.ID <= 0 {
		return false, errors.New(e.GetMsg(e.ERROR_NOT_EXIST_ORDER))
	}
	if modelO.Status >= models.ORDER_PAID {
		return true, nil
	}
	return false, nil
}

func (o *Order) Create() (*models.OrderVo, error) {
	carts, err := models.SelectCartsByUserID(o.UserId)
	if err != nil {
		log.Errorf("查询用户:%v的购物车失败 err:%v", o.UserId, err)
		return nil, err
	}
	if carts == nil {
		log.Error("购物车列表为空")
		return nil, errors.New("购物车列表为空")
	}
	orderItems, err := getCartOrderItem(o.UserId, carts)
	if err != nil {
		log.Errorf("查询购物车订单item失败 userID:%v err:%v", o.UserId, err)
		return nil, err
	}
	if orderItems == nil {
		log.Error("订单items为空")
		return nil, errors.New("订单items为空")
	}
	payment := getOrderTotalPrice(orderItems)
	order, err := assembleOrder(o.UserId, o.ShippingId, payment)
	if err != nil {
		log.Errorf("构建order实体失败 userID:%v shippingID:%v payment:%v", o.UserId, o.ShippingId, payment)
		return nil, err
	}
	for _, orderItem := range orderItems {
		orderItem.OrderNo = order.OrderNo
	}
	if err = models.BatchInsert(orderItems); err != nil {
		log.Errorf("批量插入orderItems失败 err:%v", err)
		return nil, err
	}
	if err = reduceProductStock(orderItems); err != nil {
		log.Errorf("批量减少库存失败 err:%v", err)
		return nil, err
	}
	if err = cleanCart(carts); err != nil {
		log.Errorf("清空购车失败 err:%v", err)
		return nil, err
	}
	orderVo := assembleOrderVo(order, orderItems)
	return orderVo, nil
}

func (o *Order) Cancel() error {
	modelOrder, err := models.SelectOrderByUserIdAndOrderNo(o.UserId, o.OrderNo)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Errorf("该用户此订单不存在，用户ID:%v 订单号:%v", o.UserId, o.OrderNo)
		} else {
			log.Errorf("获取订单失败，用户ID:%v 订单号:%v", o.UserId, o.OrderNo)
		}
		return err
	}
	if modelOrder.Status != models.ORDER_NOT_PAID {
		log.Errorf("订单已支付，无法取消，用户ID:%v 订单号:%v", o.UserId, o.OrderNo)
		return errors.New("订单已支付，无法取消")
	}
	err = models.UpdateOrder(o.OrderNo, map[string]interface{}{
		"status": models.ORDER_PAY_CANCEL,
	})
	if err != nil {
		log.Errorf("取消订单失败，用户ID:%v 订单号:%v", o.UserId, o.OrderNo)
	}
	return err
}

func (o *Order) GetOrderDetail() (*models.OrderVo, error) {
	order, err := models.SelectOrderByUserIdAndOrderNo(o.UserId, o.OrderNo)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Errorf("该用户此订单不存在，用户ID:%v 订单号:%v", o.UserId, o.OrderNo)
		} else {
			log.Errorf("获取订单失败，用户ID:%v 订单号:%v", o.UserId, o.OrderNo)
		}
		return nil, err
	}
	orderItems, err := models.SelectOrderItemsWithOrderNoAndUserId(o.OrderNo, o.UserId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Errorf("orderItems不存在，用户ID:%v 订单号:%v", o.UserId, o.OrderNo)
		} else {
			log.Errorf("获取orderItems失败，用户ID:%v 订单号:%v", o.UserId, o.OrderNo)
		}
		return nil, err
	}
	orderVo := assembleOrderVo(order, orderItems)
	return orderVo, nil
}

func (o *Order) GetOrderCartProducts() (*models.OrderProductVo, error) {
	carts, err := models.SelectCartsByUserID(o.UserId)
	if err != nil {
		log.Infof("查询用户:%v的购物车失败 err:%v", o.UserId, err)
		return nil, err
	}
	if carts == nil {
		log.Info("购物车列表为空")
		return nil, gorm.ErrRecordNotFound
	}
	orderItems, err := getCartOrderItem(o.UserId, carts)
	if err != nil {
		log.Infof("查询购物车订单item失败 userID:%v err:%v", o.UserId, err)
		return nil, err
	}
	if orderItems == nil {
		log.Info("订单items为空")
		return nil, gorm.ErrRecordNotFound
	}
	totalPrice := getOrderTotalPrice(orderItems)
	productVo := models.OrderProductVo{
		OrderItemList: orderItems,
		TotalPrice:    totalPrice,
		ImageHost:     "",
	}
	return &productVo, nil
}

func (o *Order) GetOrderList() ([]*models.Order, error) {
	orders, err := models.SelectOrderListByUserId(o.UserId, o.PageNum, o.PageSize)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *Order) ManageGetOrderList() ([]*models.Order, error) {
	orders, err := models.SelectOrderList(o.PageNum, o.PageSize)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *Order) ManageGetOrderDetail() (*models.OrderVo, error) {
	order, err := models.SelectOrderByOrderNo(o.OrderNo)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Errorf("此订单不存在，订单号:%v", o.OrderNo)
		} else {
			log.Errorf("获取订单失败，订单号:%v", o.OrderNo)
		}
		return nil, err
	}
	orderItems, err := models.SelectOrderItemsWithOrderNo(o.OrderNo)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Errorf("此orderItems不存在，订单号:%v", o.OrderNo)
		} else {
			log.Errorf("获取orderItems失败，订单号:%v", o.OrderNo)
		}
		return nil, err
	}
	orderVo := assembleOrderVo(order, orderItems)
	return orderVo, nil
}

func (o *Order) ManageOrderSearch() ([]*models.Order, error) {
	// TODO
	return nil, nil
}

func (o *Order) SendGoods() error {
	sendTime := time.Now()
	return models.UpdateOrder(o.OrderNo, map[string]interface{}{
		"status":    models.ORDER_SEND,
		"send_time": &sendTime,
	})
}

func getCartOrderItem(userId int, carts []*models.Cart) ([]*models.OrderItem, error) {
	var orderItems []*models.OrderItem
	for _, cart := range carts {
		orderItem := models.OrderItem{}
		product, err := models.SelectProductById(cart.ProductId)
		if err != nil {
			return nil, err
		}
		if product.Status != models.OnSale {
			return nil, errors.New("产品" + product.Name + "不是在线销售状态")
		}
		if cart.Quantity > product.Stock {
			return nil, errors.New("产品" + product.Name + "库存不足")
		}
		orderItem.UserId = userId
		orderItem.ProductId = product.ID
		orderItem.ProductName = product.Name
		orderItem.ProductImage = product.MainImage
		orderItem.CurrentUnitPrice = product.Price
		orderItem.Quantity = cart.Quantity
		orderItem.TotalPrice = product.Price * cart.Quantity
		orderItems = append(orderItems, &orderItem)
	}
	return orderItems, nil
}

func getOrderTotalPrice(orderItems []*models.OrderItem) int {
	totalPrice := 0
	for _, orderItem := range orderItems {
		totalPrice += orderItem.TotalPrice
	}
	return totalPrice
}

func assembleOrderVo(order *models.Order, orderItems []*models.OrderItem) *models.OrderVo {
	orderVo := models.OrderVo{
		OrderNo:       order.OrderNo,
		Payment:       order.Payment,
		PaymentType:   order.PaymentType,
		PaymentDesc:   "在线支付",
		Postage:       order.Postage,
		Status:        order.Status,
		StatusDesc:    "",
		ShippingId:    order.ShippingId,
		PaymentTime:   order.PaymentTime,
		SendTime:      order.SendTime,
		EndTime:       order.EndTime,
		CloseTime:     order.CloseTime,
		CreateTime:    &order.CreatedAt,
		ImageHost:     "",
		OrderItemList: orderItems,
	}
	shipping := models.Shipping{
		Model:  models.Model{ID: order.ShippingId},
		UserId: order.UserId,
	}
	err := shipping.Get()
	if err != nil || shipping.ID > 0 {
		orderVo.Shipping = &shipping
	}
	return &orderVo
}

func assembleOrder(userId, shippingId, payment int) (*models.Order, error) {
	model := models.Order{
		OrderNo:     generateOrderNo(),
		Status:      models.ORDER_NOT_PAID,
		Postage:     0,
		PaymentType: models.PaymentTypeOnlinePayment,
		Payment:     payment,
		UserId:      userId,
		ShippingId:  shippingId,
	}
	err := model.Create()
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func generateOrderNo() int64 {
	return time.Now().UnixNano() + int64(rand.Intn(100))
}

func reduceProductStock(orderItems []*models.OrderItem) error {
	var err error
	for _, orderItem := range orderItems {
		p, err := models.SelectProductById(orderItem.ProductId)
		if err != nil {
			return err
		}
		p.Stock -= orderItem.Quantity
		err = p.Update()
	}
	return err
}

func cleanCart(carts [] *models.Cart) error {
	var err error
	for _, cart := range carts {
		err = models.DeleteCart(cart.ID)
	}
	return err
}
