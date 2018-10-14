package models

const (
	PAY_PLATFORM_ALIPAY = 1
	PAY_PLATFORM_WeChat = 2
)

const (
	TRADE_STATUS_WAIT_BUYER_PAY = "WAIT_BUYER_PAY"
	TRADE_STATUS_TRADE_SUCCESS  = "TRADE_SUCCESS"
	RESPONSE_SUCCESS = "success"
	RESPONSE_FAILED  = "failed"
)

type PayInfo struct {
	Model

	UserId         int    `json:"user_id"`         // 用户ID
	OrderNo        int    `json:"order_no"`        // 订单号
	PayPlatform    int    `json:"pay_platform"`    // 支付平台 1-支付宝 2-微信
	PlatformNumber string `json:"platform_number"` // 支付平台返回的流水号
	PlatformStatus string `json:"platform_status"` // 支付状态
}

func SelectPayInfoByUserIdAndOrderNo(userId, orderNo int) (*PayInfo, error) {
	var payInfo PayInfo
	err := db.First(&payInfo, "user_id = ? and order_no = ?", userId, orderNo).Error
	if err != nil {
		return nil, err
	}
	return &payInfo, nil
}

func UpdatePayInfo(id int, data map[string]interface{}) error {
	return db.Model(&PayInfo{}).Where("id = ?", id).Updates(data).Error
}

func (p *PayInfo)Create() error {
	return db.Create(p).Error
}
