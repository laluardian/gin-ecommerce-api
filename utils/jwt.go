package utils

import (
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/laluardian/gin-ecommerce-api/models"
	"github.com/rs/xid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token is expired")
)

const JwtPayloadKey = "jwt_payload"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type JwtPayload struct {
	Sub      xid.ID    `json:"sub"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	Iat      time.Time `json:"iat"`
	Exp      time.Time `json:"exp"`
}

func newJwtPayload(user *models.User) *JwtPayload {
	role := RoleUser
	if user.IsAdmin {
		role = RoleAdmin
	}

	return &JwtPayload{
		Sub:      user.ID,
		Username: user.Username,
		Role:     role,
		Iat:      time.Now(),
		Exp:      time.Now().Add(time.Hour * 24),
	}
}

func (p *JwtPayload) Valid() error {
	if time.Now().After(p.Exp) {
		return ErrExpiredToken
	}

	return nil
}

func GenerateToken(user *models.User) (string, error) {
	payload := newJwtPayload(user)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyToken(token string) (*JwtPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &JwtPayload{}, keyFunc)
	if err != nil {
		valErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(valErr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*JwtPayload)
	if !jwtToken.Valid || !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

func CheckUserId(c *gin.Context, userId xid.ID) *JwtPayload {
	authPayload := c.MustGet(JwtPayloadKey).(*JwtPayload)
	if authPayload.Sub != userId {
		return nil
	}

	return authPayload
}

func CheckUserRole(c *gin.Context) *JwtPayload {
	authPayload := c.MustGet(JwtPayloadKey).(*JwtPayload)
	if authPayload.Role != RoleAdmin {
		return nil
	}

	return authPayload
}
