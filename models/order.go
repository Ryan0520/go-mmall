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

	OrderNo     int64      `json:"order_no"`
	UserId      int        `json:"user_id"`
	ShippingId  int        `json:"shipping_id"`
	Payment     int        `json:"payment"` // 实际付款金额，单位分
	PaymentType int        `json:"payment_type"`
	Postage     int        `json:"postage"`                 // 运费，单位分
	Status      int        `json:"status"`                  // 订单状态
	PaymentTime *time.Time `json:"payment_time, omitempty"` // 支付时间
	SendTime    *time.Time `json:"send_time, omitempty"`    // 发货时间
	EndTime     *time.Time `json:"end_time, omitempty"`     // 交易完成时间
	CloseTime   *time.Time `json:"close_time, omitempty"`   // 交易关闭时间
}

type OrderVo struct {
	OrderNo       int64        `json:"order_no"`
	Payment       int          `json:"payment"`
	PaymentType   int          `json:"payment_type"`
	PaymentDesc   string       `json:"payment_desc"`
	Postage       int          `json:"postage"`
	Status        int          `json:"status"`
	StatusDesc    string       `json:"status_desc"`
	CreateTime    *time.Time   `json:"create_time"`
	PaymentTime   *time.Time   `json:"payment_time, omitempty"` // 支付时间
	SendTime      *time.Time   `json:"send_time, omitempty"`    // 发货时间
	EndTime       *time.Time   `json:"end_time, omitempty"`     // 交易完成时间
	CloseTime     *time.Time   `json:"close_time, omitempty"`   // 交易关闭时间
	OrderItemList []*OrderItem `json:"order_item_list"`
	ShippingId    int          `json:"shipping_id"`
	Shipping      *Shipping    `json:"shipping"`
	ImageHost     string       `json:"image_host"`
}

type OrderProductVo struct {
	OrderItemList []*OrderItem `json:"order_item_list"`
	TotalPrice    int          `json:"total_price"`
	ImageHost     string       `json:"image_host"`
}

func SelectOrderByUserIdAndOrderNo(userId int, orderNo int64) (*Order, error) {
	var order Order
	err := db.First(&order, "user_id = ? and order_no = ?", userId, orderNo).Error
	if err != nil {
		return nil, err
	}
	return &order, err
}

func SelectOrderByOrderNo(orderNo int64) (*Order, error) {
	var order Order
	err := db.First(&order, "order_no = ?", orderNo).Error
	if err != nil {
		return nil, err
	}
	return &order, err
}

func SelectOrderListByOrderNo(orderNo int64, pageNum, pageSize int) ([]*Order, error) {
	var orders []*Order
	err := db.Find(&orders, "order_no = ?", orderNo).Offset(pageNum).Limit(pageSize).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func SelectOrderListByUserId(userId, pageNum, pageSize int) ([]*Order, error) {
	var orders []*Order
	err := db.Find(&orders, "user_id = ?", userId).Offset(pageNum).Limit(pageSize).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func SelectOrderList(pageNum, pageSize int) ([]*Order, error) {
	var orders []*Order
	err := db.Find(&orders).Offset(pageNum).Limit(pageSize).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func UpdateOrder(orderNo int64, data map[string]interface{}) error {
	err := db.Model(&Order{}).Where("order_no = ?", orderNo).Updates(data).Error
	return err
}

func (o *Order) Create() error {
	return db.Create(o).Error
}
