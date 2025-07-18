package models

import (
	"html"
	"time"
	"net/url"

	"github.com/satori/uuid"
)

// easyjson:json
type Advertisement struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}

func (a *Advertisement) Sanitize() {
	a.Title = html.EscapeString(a.Title)
	a.Description = html.EscapeString(a.Description)
	a.ImageURL = url.QueryEscape(a.ImageURL)
}
