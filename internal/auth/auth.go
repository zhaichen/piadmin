package auth

import (
	"crypto/sha256"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ptmind/piadmin/internal/config"
)

var ErrInvalidPassword = errors.New("invalid password")
var ErrInvalidToken = errors.New("invalid token")

type Auth struct {
	cfg    config.AuthConfig
	secret []byte
}

func New(cfg config.AuthConfig) *Auth {
	h := sha256.Sum256([]byte(cfg.Password))
	return &Auth{
		cfg:    cfg,
		secret: h[:],
	}
}

func (a *Auth) Login(password string) (string, error) {
	if password != a.cfg.Password {
		return "", ErrInvalidPassword
	}
	return a.generateToken()
}

func (a *Auth) ValidateToken(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return a.secret, nil
	})
	if err != nil || !token.Valid {
		return ErrInvalidToken
	}
	return nil
}

func (a *Auth) generateToken() (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(a.cfg.TokenTTL) * time.Second)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secret)
}
