package validation

import (
	"errors"
	"strings"
	"testing"

	"github.com/K1tten2005/go_vk_intern/internal/models"
)

func TestHashPasswordAndCheckPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"simple password", "Password123!"},
		{"complex password", "P@ssw0rd!VeryL0ngAndSecure"},
		{"short password", "P@ss1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			salt := []byte("randsalt")
			hashed := HashPassword(salt, tt.password)

			if !CheckPassword(hashed, tt.password) {
				t.Errorf("CheckPassword() failed for password: %s", tt.password)
			}

			if CheckPassword(hashed, "wrongpassword") {
				t.Error("CheckPassword() should return false for wrong password")
			}
		})
	}
}

func TestValidTextContent(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		maxLen int
		want   bool
	}{
		{"empty string", "", 10, false},
		{"valid text", "Valid text 123", 20, true},
		{"invalid chars", "Invalid@text", 20, false},
		{"russian chars", "Привет мир", 20, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidTextContent(tt.text, tt.maxLen); got != tt.want {
				t.Errorf("ValidTextContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidTitle(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  bool
	}{
		{"empty title", "", false},
		{"too long", strings.Repeat("a", maxTitleLength+1), false},
		{"valid title", "Valid Title 123", true},
		{"invalid chars", "Invalid@Title", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidTitle(tt.title); got != tt.want {
				t.Errorf("ValidTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		want        bool
	}{
		{"empty description", "", false},
		{"too long", strings.Repeat("a", maxDescriptionLength+1), false},
		{"valid description", "Valid Description 123, with some symbols: #*,.", true},
		{"invalid chars", "Invalid@Description", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidDescription(tt.description); got != tt.want {
				t.Errorf("ValidDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidImageURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{"empty url", "", false},
		{"too long", "http://example.com/" + strings.Repeat("a", maxImageURLLength), false},
		{"invalid url", "not a url", false},
		{"no scheme", "example.com/image.jpg", false},
		{"no host", "http:///image.jpg", false},
		{"valid jpg", "http://example.com/image.jpg", true},
		{"valid jpeg", "https://example.com/image.jpeg", true},
		{"valid png", "http://example.com/image.png", true},
		{"valid webp", "https://example.com/image.webp", true},
		{"invalid extension", "http://example.com/image.gif", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidImageURL(tt.url); got != tt.want {
				t.Errorf("ValidImageURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidPrice(t *testing.T) {
	tests := []struct {
		name  string
		price int
		want  bool
	}{
		{"negative", -1, false},
		{"zero", 0, true},
		{"normal", 1000, true},
		{"max", MaxPrice, true},
		{"over max", MaxPrice + 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidPrice(tt.price); got != tt.want {
				t.Errorf("ValidPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateAd(t *testing.T) {
	validAd := models.Ad{
		Title:       "Valid Title",
		Description: "Valid description with some text",
		ImageURL:    "http://example.com/image.jpg",
		Price:       1000,
	}

	tests := []struct {
		name    string
		modify  func(*models.Ad)
		wantErr error
	}{
		{"valid", func(a *models.Ad) {}, nil},
		{"invalid title", func(a *models.Ad) { a.Title = "" }, errors.New("invalid title")},
		{"invalid description", func(a *models.Ad) { a.Description = "" }, errors.New("invalid description")},
		{"invalid image url", func(a *models.Ad) { a.ImageURL = "invalid" }, errors.New("invalid image url")},
		{"invalid price", func(a *models.Ad) { a.Price = -1 }, errors.New("invalid price")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ad := validAd
			tt.modify(&ad)

			err := ValidateAd(ad)
			if (err == nil) != (tt.wantErr == nil) {
				t.Errorf("ValidateAd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("ValidateAd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidLogin(t *testing.T) {
	tests := []struct {
		name  string
		login string
		want  bool
	}{
		{"too short", "ab", false},
		{"too long", strings.Repeat("a", maxLoginLength+1), false},
		{"valid", "user_123", true},
		{"invalid chars", "user@name", false},
		{"min length", "abc", true},
		{"max length", strings.Repeat("a", maxLoginLength), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidLogin(tt.login); got != tt.want {
				t.Errorf("ValidLogin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"too short", "P@ss1", false},
		{"too long", strings.Repeat("P", maxPassLength+1) + "@1", false},
		{"no upper", "password@1", false},
		{"no lower", "PASSWORD@1", false},
		{"no digit", "Password@", false},
		{"no special", "Password1", false},
		{"valid", "P@ssword123", true},
		{"max length", "P@ss" + strings.Repeat("w", maxPassLength-5) + "1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidPassword(tt.password); got != tt.want {
				t.Errorf("ValidPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
