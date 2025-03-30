package model

const (
	PATH_USER_ROOT         string = "user"
	PATH_USER_TABLE_CREATE string = "user/register"
)

type User struct {
	UserID    int    `json:"userid"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}

func (u *User) Mapping() map[string]any {
	if u == nil {
		return nil
	}
	mp := make(map[string]any, 3)
	if u.UserID > 0 {
		mp["id"] = u.UserID
	}
	if u.Name != "" {
		mp["name"] = u.Name
	}
	if u.Telephone != "" {
		mp["tel"] = u.Telephone
	}
	return mp
}
