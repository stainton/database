package model

const (
	PATH_REWARD_ROOT         string = "reward"
	PATH_REWARD_TABLE_CREATE string = "reward/register"
)

type Reward struct {
	ProductId int64  `json:"productid"`
	Date      string `json:"date"`
	Type      string `json:"type"`
}

// productid
// type
// date

func (r *Reward) Mapping() map[string]any {
	if r == nil {
		return nil
	}
	mp := make(map[string]any, 3)
	if r.ProductId <= 0 {
		mp["productid"] = r.ProductId
	}
	if r.Date != "" {
		mp["date"] = r.Date
	}
	if r.Type != "" {
		mp["type"] = r.Type
	}
	return mp
}
