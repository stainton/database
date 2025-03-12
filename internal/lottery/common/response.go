package common

type ResponseTemplate struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type UserPay struct {
	Userid int `json:"userid"`
	Pay    int `json:"pay"`
}

type UserPayList struct {
	Productid int        `json:"productid"`
	Users     []*UserPay `json:"users"`
}
