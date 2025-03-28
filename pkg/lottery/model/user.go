package model

type User struct {
	UserID    int    `json:"userid"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}

const (
	PATH_USER_ROOT string = "user"
)
