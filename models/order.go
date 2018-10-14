package models

import "time"

// 订单状态
const (
	// 已取消
	ORDER_PAY_CANCEL = 0
	// 未支付
	ORDER_NOT_PAID = 10
	// 已支付
	ORDER_PAID = 20
	// 已发货
	ORDER_SEND = 40
	// 交易成功
	ORDER_SUCCESS = 50
	// 交易关闭
	ORDER_CLOSE = 60
)

// 支付类型
const (
	// 在线支付
	PaymentTypeOnlinePayment = 1
)

type Order struct {
	Model

	OrderNo     int       `json:"order_no"`
	UserId      int       `json:"user_id"`
	ShippingId  int       `json:"shipping_id"`
	Payment     int       `json:"payment"` // 实际付款金额，单位分
	PaymentType int       `json:"payment_type"`
	Postage     int       `json:"postage"`                 // 运费，单位分
	Status      int       `json:"status"`                  // 订单状态
	PaymentTime time.Time `json:"payment_time, omitempty"` // 支付时间
	SendTime    time.Time `json:"send_time, omitempty"`    // 发货时间
	EndTime     time.Time `json:"end_time, omitempty"`     // 交易完成时间
	CloseTime   time.Time `json:"close_time, omitempty"`   // 交易关闭时间
}

func SelectOrderByUserIdAndOrderNo(userId, orderNo int) (*Order, error) {
	var order Order
	err := db.First(&order, "user_id = ? and order_no = ?", userId, orderNo).Error
	if err != nil {
		return nil, err
	}
	return &order, err
}

func SelectOrderByOrderNo(orderNo int) (*Order, error) {
	var order Order
	err := db.First(&order, "order_no = ?", orderNo).Error
	if err != nil {
		return nil, err
	}
	return &order, err
}

func UpdateOrder(orderNo int, data map[string]interface{}) error {
	err := db.Model(&Order{}).Where("order_no = ?", orderNo).Updates(data).Error
	return err
}