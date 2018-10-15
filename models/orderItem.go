package models

type OrderItem struct {
	Model

	UserId int `json:"user_id"`
	OrderNo int64 `json:"order_no"`
	ProductId int `json:"product_id"`
	ProductName string `json:"product_name"`
	ProductImage string `json:"product_image"`
	CurrentUnitPrice int `json:"current_unit_price"`
	Quantity int `json:"quantity"`
	TotalPrice int `json:"total_price"`
}

func SelectOrderItemsWithOrderNoAndUserId(orderNo int64, userId int) ([]*OrderItem, error) {
	var orderItems []*OrderItem
	err := db.Find(&orderItems, "order_no = ? and user_id = ?", orderNo, userId).Error
	if err != nil {
		return nil, err
	}
	return orderItems, nil
}

func SelectOrderItemsWithOrderNo(orderNo int64) ([]*OrderItem, error) {
	var orderItems []*OrderItem
	err := db.Find(&orderItems, "order_no = ?", orderNo).Error
	if err != nil {
		return nil, err
	}
	return orderItems, nil
}

func BatchInsert(orderItems []*OrderItem) error {
	tx := db.Begin()
	var err error
	for _, orderItem := range orderItems {
		err = tx.Create(&orderItem).Error
	}
	if err != nil {
		tx.Rollback()
	}
 	err = tx.Commit().Error
	return err

	//values := ""
	//for _, orderItem := range orderItems {
	//	values = values + "(" + strconv.Itoa(orderItem.ID) + "," + fmt.Sprintf("%v", orderItem.OrderNo) + strconv.Itoa(orderItem.UserId) +
	//		strconv.Itoa(orderItem.ProductId) + orderItem.ProductName + orderItem.ProductImage + strconv.Itoa(orderItem.CurrentUnitPrice) +
	//		strconv.Itoa(orderItem.Quantity) + strconv.Itoa(orderItem.TotalPrice) + ")"
	//}
	//fmt.Println(values)
	//
	//sqlRaw := fmt.Sprintf("insert into mmall_order_item (id, order_no, user_id, product_id, product_name, product_image, " +
	//	"current_unit_price, quantity, total_price) values %s", values)
	//return db.Exec(sqlRaw).Error
}
