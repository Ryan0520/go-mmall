package portal

import (
	"github.com/Ryan0520/go-mmall/models"
	"github.com/Ryan0520/go-mmall/pkg/e"
	"github.com/Ryan0520/go-mmall/pkg/setting"
	ali "github.com/Ryan0520/go-mmall/service/alipay"
	"github.com/kataras/iris/core/errors"
	log "github.com/sirupsen/logrus"
	"github.com/smartwalle/alipay"
	"strconv"
)

type Order struct {
	OrderNo int    `json:"order_no"`
	UserId  int    `json:"user_id"`
	PayUrl  string `json:"pay_url"`
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

	log.Infoln("========== TradePagePay ==========")
	var p = alipay.AliPayTradePagePay{}
	p.NotifyURL = setting.AlipaySetting.NotifyUrl
	p.ReturnURL = setting.AlipaySetting.ReturnUrl
	p.Subject = "订单" + strconv.Itoa(o.OrderNo) + "购买商品共" + strconv.Itoa(modelO.Payment/100) + "元"
	p.OutTradeNo = strconv.Itoa(o.OrderNo)
	p.TotalAmount = strconv.Itoa(modelO.Payment/100)
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