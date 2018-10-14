package alipay_test

import (
	ali "github.com/Ryan0520/go-mmall/service/alipay"
	"github.com/smartwalle/alipay"
	"testing"
)

func TestAliPay_TradePagePay(t *testing.T) {
	t.Log("========== TradePagePay ==========")
	var p = alipay.AliPayTradePagePay{}
	p.NotifyURL = "http://bzzyq2.natappfree.cc/alipay/notify"
	p.ReturnURL = "http://bzzyq2.natappfree.cc/alipay/return"
	p.Subject = "测试支付"
	p.OutTradeNo = "trade_no_2017062301223334112"
	p.TotalAmount = "2.00"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	url, err := ali.Client.TradePagePay(p)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(url)
}
