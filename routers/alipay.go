package routers

import (
	"github.com/Ryan0520/go-mmall/models"
	"github.com/Ryan0520/go-mmall/service/alipay"
	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func NotifyHandle(c *gin.Context) {
	var notification, _ = alipay.Client.GetTradeNotification(c.Request)
	if notification != nil {
		log.Println("支付成功")
		orderNo := com.StrTo(notification.OutTradeNo).MustInt()
		order, err := models.SelectOrderByOrderNo(orderNo)
		if err != nil || order.ID <= 0{
			log.Errorf("获取订单失败，订单号: %d", orderNo)
			c.String(http.StatusOK, models.RESPONSE_FAILED)
			return
		}
		if order.Status >= models.ORDER_PAID {
			log.Infof("订单号: %d 是已支付状态，支付宝重复回调", orderNo)
			c.String(http.StatusOK, models.RESPONSE_FAILED)
			return
		}
		order.Status = models.ORDER_PAID
		order.PaymentTime = time.Now()
		err = models.UpdateOrder(order.OrderNo, map[string]interface{}{
			"payment_time": order.PaymentTime,
			"status": order.Status,
		})
		if err != nil {
			log.Errorf("更新订单状态失败，订单号: %d", orderNo)
			c.String(http.StatusOK, models.RESPONSE_FAILED)
			return
		}

		p := models.PayInfo{
			UserId: order.UserId,
			OrderNo: orderNo,
			PayPlatform: models.PAY_PLATFORM_ALIPAY,
			PlatformNumber: notification.TradeNo,
			PlatformStatus: notification.TradeStatus,
		}
		err = p.Create()
		if err != nil {
			log.Errorf("创建订单支付信息失败 err: %v", err)
			c.String(http.StatusOK, models.RESPONSE_FAILED)
			return
		}
		log.Info("创建订单支付信息成功，回调支付宝success")
		c.String(http.StatusOK, models.RESPONSE_SUCCESS)
	} else {
		log.Println("支付失败")
		c.String(http.StatusOK, models.RESPONSE_FAILED)
	}
}

func ReturnHandle(c *gin.Context) {
	c.Request.ParseForm()
	ok, err := alipay.Client.VerifySign(c.Request.Form)
	log.Info(ok, err)
	status := "success"
	if !ok {
		status = "failure"
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"name": "Ryan",
		"status": status,
	})
}
