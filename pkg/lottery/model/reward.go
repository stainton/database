package model

type Reward struct {
	ProductId int64  `json:"productid"`
	Date      string `json:"date"`
	Type      string `json:"type"`
}
