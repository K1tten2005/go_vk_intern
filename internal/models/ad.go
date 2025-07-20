package models

import (
	"html"
	"time"

	"github.com/satori/uuid"
)

// easyjson:json
type Ad struct {
	Id          uuid.UUID
	UserId      uuid.UUID
	Title       string
	Description string
	Price       int
	ImageURL    string
	CreatedAt   time.Time
	AuthorLogin string
}

// easyjson:json
type AdReq struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
}

// easyjson:json
type AdResp struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	AuthorLogin string    `json:"author_login"`
	IsOwner     bool      `json:"is_owner,omitempty"`
}

//easyjson:json
type AdRespList []AdResp

type Filter struct {
	Page     int
	Limit    int
	SortBy   string
	Order    string
	PriceMin int
	PriceMax int
	UserId   uuid.UUID
}

func (a *Ad) Sanitize() {
	a.Title = html.EscapeString(a.Title)
	a.Description = html.EscapeString(a.Description)
}
