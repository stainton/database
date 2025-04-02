package model

const (
	PATH_ORDER_ROOT         string = "order"
	PATH_ORDER_TABLE_CREATE string = "order/register"
)

type Order struct {
	OrderId   int64  `json:"orderid"`
	UserId    int64  `json:"userid"`
	ProductId int64  `json:"productid"`
	Price     int64  `json:"price"`
	Date      string `json:"date"`
	Type      string `json:"type"`
}

// orderid,userid,productid,price,create_time,type

func (r *Order) Mapping() map[string]any {
	if r == nil {
		return nil
	}
	mp := make(map[string]any, 3)
	if r.OrderId > 0 {
		mp["orderid"] = r.OrderId
	}
	if r.UserId > 0 {
		mp["userid"] = r.UserId
	}
	if r.ProductId > 0 {
		mp["productid"] = r.ProductId
	}
	if r.Price > 0 {
		mp["price"] = r.Price
	}
	if r.Date != "" {
		mp["date"] = r.Date
	}
	if r.Type != "" {
		mp["type"] = r.Type
	}
	return mp
}
