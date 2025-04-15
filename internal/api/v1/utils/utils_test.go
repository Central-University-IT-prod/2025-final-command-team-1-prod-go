package utils

import (
	"net/http"
	"testing"
	"time"

	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestExcludeUserCredentials(t *testing.T) {
	userDto := dto.UserDto{
		Username:  "test",
		Email:     "email@email.com",
		CreatedAt: "2021-01-01T00:00:00Z",
		UpdatedAt: "2021-01-01T00:00:00Z",
		Password:  "password",
	}

	got := ExcludeUserCredentials(&userDto)
	want := dto.GetUserDto{
		Username:  "test",
		Email:     "email@email.com",
		CreatedAt: "2021-01-01T00:00:00Z",
		UpdatedAt: "2021-01-01T00:00:00Z",
	}

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestHashPassword(t *testing.T) {
	password := "password"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("got error %v", err)
	}

	if hashedPassword == password {
		t.Errorf("got %v want %v", hashedPassword, password)
	}
}

func TestUpdateUserTimestamps(t *testing.T) {
	updateUserDto := dto.UpdateUserDto{
		Username:  "test",
		Password:  "password",
		UpdatedAt: "2021-01-01T00:00:00Z",
	}

	UpdateUserTimestamps(&updateUserDto)

	if updateUserDto.UpdatedAt == "2021-01-01T00:00:00Z" {
		t.Errorf("got %v want different value", updateUserDto.UpdatedAt)
	}
}

func TestExtractTokenFromHeaders(t *testing.T) {
	// test case when Authorization header is not provided
	c := &gin.Context{}
	c.Request = &http.Request{}
	c.Request.Header = make(http.Header)
	_, err := ExtractTokenFromHeaders(c)

	if err == nil {
		t.Errorf("got %v want error", err)
	}

	// test case when token is provided
	c.Request.Header.Set("Authorization", "Bearer token")
	token, err := ExtractTokenFromHeaders(c)
	if err != nil {
		t.Errorf("got error %v", err)
	}

	if *token != "token" {
		t.Errorf("got %v want token", *token)
	}
}

func TestValidateTokenSignature(t *testing.T) {
	// test case when token is invalid
	err := ValidateTokenSignature("invalid token")
	if err == nil {
		t.Errorf("got %v want error", err)
	}

	// test case when token is valid
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    "email@email.com",
		"username": "username",
		"exp":      time.Now().UTC().Add(time.Hour * 24 * 90).Unix(),
		"iat":      time.Now().UTC().Unix(),
	})
	tokenString, _ := token.SignedString([]byte(config.Config.JWTSecret))

	err = ValidateTokenSignature(tokenString)
	if err != nil {
		t.Errorf("got error %v", err)
	}
}

func TestExtractPayloadFromJWT(t *testing.T) {
	// test case when token is invalid
	_, err := ExtractPayloadFromJWT("invalid token")
	if err == nil {
		t.Errorf("got %v want error", err)
	}

	// test case when token is valid
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    "email@email.com",
		"username": "username",
		"exp":      time.Now().UTC().Add(time.Hour * 24 * 90).Unix(),
		"iat":      time.Now().UTC().Unix(),
	})
	tokenString, _ := token.SignedString([]byte(config.Config.JWTSecret))

	claims, err := ExtractPayloadFromJWT(tokenString)
	if err != nil {
		t.Errorf("got error %v", err)
	}

	if claims["email"] != "email@email.com" {
		t.Errorf("got %v want email", claims["email"])
	}

	if claims["username"] != "username" {
		t.Errorf("got %v want username", claims["username"])
	}
}
