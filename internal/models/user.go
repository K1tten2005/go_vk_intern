package models

import (
	"html"

	"github.com/satori/uuid"
)

// easyjson:json
type User struct {
	Id           uuid.UUID `json:"id"`
	Login        string    `json:"login"`
	PasswordHash []byte    `json:"-"`
}

// easyjson:json
type UserReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// easyjson:json
type UserResp struct {
	Id    uuid.UUID `json:"id"`
	Login string    `json:"login"`
	Token string    `json:"token"`
}

func (u *User) Sanitize() {
	u.Login = html.EscapeString(u.Login)
}

func (u *UserReq) Sanitize() {
	u.Login = html.EscapeString(u.Login)
	u.Password = html.EscapeString(u.Password)
}

func (u *UserResp) Sanitize() {
	u.Login = html.EscapeString(u.Login)
}
