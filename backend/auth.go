package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// jwtSecret signs admin session tokens. Override with the JWT_SECRET
// env var in production — this fallback is fine for local dev only.
var jwtSecret = []byte(getEnvOr("JWT_SECRET", "dev-only-change-me"))

func getEnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

func login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}

	var user User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	token, err := issueToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue token"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{Token: token, Username: user.Username})
}

func issueToken(user User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"usr": user.Username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// authRequired protects mutation routes (create/update/delete).
// Public GET routes stay open so the storefront can read data without logging in.
func authRequired(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
		return
	}
	tokenString := strings.TrimPrefix(header, "Bearer ")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	c.Next()
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// changePassword lets the logged-in admin rotate their own password.
func changePassword(c *gin.Context) {
	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "current_password and new_password (min 8 chars) required"})
		return
	}

	// Single-admin setup: there's one user row. For multi-user, pull
	// the user id out of the JWT claims (c.Get after middleware sets it) instead.
	var user User
	if err := db.First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no admin account found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "current password is incorrect"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user.PasswordHash = string(hash)
	db.Save(&user)
	c.JSON(http.StatusOK, gin.H{"updated": true})
}
