package validation

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/K1tten2005/go_vk_intern/internal/models"
	"golang.org/x/crypto/argon2"
)

const (
	maxLoginLength       = 20
	minLoginLength       = 3
	minPassLength        = 8
	maxPassLength        = 25
	maxTitleLength       = 100
	maxDescriptionLength = 700
	maxImageURLLength    = 300
	MaxPrice             = 100000000
	maxImageSizeBytes    = 10 * 1024 * 1024
	allowedChars         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"

	allowedSymbolsForText = "абвгдеёжзийклмнопрстуфхцчшщъыьэюя" +
		"АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ" +
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 -_#*,./!?: #*,. "
)

func HashPassword(salt []byte, plainPassword string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), salt, 1, 64*1024, 4, 32)
	return append(salt, hashedPass...)
}

func CheckPassword(passHash []byte, plainPassword string) bool {
	salt := make([]byte, 8)
	copy(salt, passHash[:8])
	userPassHash := HashPassword(salt, plainPassword)
	return bytes.Equal(userPassHash, passHash)
}

func ValidTextContent(s string, maxLen int) bool {
	if len(s) == 0 || len(s) > maxLen {
		return false
	}
	for _, r := range s {
		if !strings.ContainsRune(allowedSymbolsForText, r) {
			return false
		}
	}
	return true
}

func ValidTitle(title string) bool {
	return ValidTextContent(title, maxTitleLength)
}

func ValidDescription(desc string) bool {
	return ValidTextContent(desc, maxDescriptionLength)
}

func ValidImageURL(link string) bool {
	if len(link) == 0 || len(link) > maxImageURLLength {
		return false
	}

	u, err := url.ParseRequestURI(link)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(link)
	if err != nil {
		fmt.Println("GET request failed:", err)
		return false
	}
	defer resp.Body.Close()

	head := make([]byte, 512)
	n, err := io.ReadFull(io.LimitReader(resp.Body, 512), head)
	if err != nil && err != io.ErrUnexpectedEOF {
		fmt.Println("Error reading image head:", err)
		return false
	}

	contentType := http.DetectContentType(head[:n])
	switch contentType {
	case "image/jpeg", "image/png", "image/webp":
		return true
	default:
		fmt.Println("Invalid image content type:", contentType)
		return false
	}
}

func ImageSizeUnderLimit(link string) bool {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(link)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer resp.Body.Close()

	limited := io.LimitReader(resp.Body, maxImageSizeBytes+1)

	n, err := io.Copy(io.Discard, limited)
	if err != nil {
		fmt.Printf("Image size for %s is %d bytes\n", link, n)
		return false
	}

	return n <= int64(maxImageSizeBytes)
}

func ValidPrice(price int) bool {
	return price >= 0 && price <= MaxPrice
}

func ValidateAd(ad models.Ad) error {
	if !ValidTitle(ad.Title) {
		return errors.New("invalid title")
	}
	if !ValidDescription(ad.Description) {
		return errors.New("invalid description")
	}
	if !ValidImageURL(ad.ImageURL) {
		return errors.New("invalid image url")
	}
	if !ImageSizeUnderLimit(ad.ImageURL) {
		return errors.New("image size exceeds limit")
	}
	if !ValidPrice(ad.Price) {
		return errors.New("invalid price")
	}
	return nil
}

func ValidLogin(login string) bool {
	if len(login) < minLoginLength || len(login) > maxLoginLength {
		return false
	}
	for _, char := range login {
		if !strings.Contains(allowedChars, string(char)) {
			return false
		}
	}
	return true
}

func ValidPassword(password string) bool {
	var up, low, digit, special bool

	if len(password) < minPassLength || len(password) > maxPassLength {
		return false
	}

	for _, char := range password {

		switch {
		case unicode.IsUpper(char):
			up = true
		case unicode.IsLower(char):
			low = true
		case unicode.IsDigit(char):
			digit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			special = true
		default:
			return false
		}
	}

	return up && low && digit && special
}
