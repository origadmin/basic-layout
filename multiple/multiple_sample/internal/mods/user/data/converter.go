package data

import (
	"basic-layout/multiple/multiple_sample/internal/mods/user/biz"
)

// toUserDO converts a User PO to a DO.
func toUserDO(p *User) *biz.User {
	if p == nil {
		return nil
	}
	return &biz.User{
		ID:       p.ID,
		Username: p.Username,
		Nickname: p.Nickname,
	}
}

// fromUserDO converts a User DO to a PO.
func fromUserDO(d *biz.User) *User {
	if d == nil {
		return nil
	}
	return &User{
		ID:       d.ID,
		Username: d.Username,
		Nickname: d.Nickname,
	}
}