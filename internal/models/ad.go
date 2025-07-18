package models

import (
	"html"
	"net/url"
	"time"

	"github.com/satori/uuid"
)

// easyjson:json
type Ad struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	IsOwner     bool      `json:"is_owner"`
}

// easyjson:json
type AdResp struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64       `json:"price"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	IsOwner     bool      `json:"is_owner"`
}

func (a *Ad) Sanitize() {
	a.Title = html.EscapeString(a.Title)
	a.Description = html.EscapeString(a.Description)
	a.ImageURL = url.QueryEscape(a.ImageURL)
}
