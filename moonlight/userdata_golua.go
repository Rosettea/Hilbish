//go:build !midnight

package moonlight

import (
	rt "github.com/arnodel/golua/runtime"
)

type UserData struct {
	ud *rt.UserData
}

func NewUserData(v interface{}, meta *Table) *UserData {
	return &UserData{
		ud: rt.NewUserData(v, meta.lt),
	}
}

func UserDataValue(u *UserData) Value {
	return rt.UserDataValue(u.ud)
}
