package model

type User struct {
	Userid   int64  `json:"userid"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
