package model

type Order struct {
	OrderId   int64  `json:"orderid"`
	UserId    int64  `json:"userid"`
	ProductId int64  `json:"productid"`
	Price     int64  `json:"price"`
	Date      string `json:"date"`
	Type      string `json:"type"`
}
