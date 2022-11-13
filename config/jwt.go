package config

import (
	"errors"
	"os"
	"time"

	"github.com/kataras/iris/v12/middleware/jwt"
)

const (
	accessTokenMaxAge  = 10 * time.Minute
	RefreshTokenMaxAge = time.Hour
)

var (
	SecretKey = []byte(os.Getenv("SECRET"))
	Signer    *jwt.Signer
    Verifier = jwt.NewVerifier(jwt.HS256, SecretKey)
)

func InitJWT() error {
    mode := os.Getenv("MODE")
    if mode == "" {
        return errors.New("Empty envvar $MODE")
    } else if mode == "development" {
		Signer = jwt.NewSigner(jwt.HS256, SecretKey, accessTokenMaxAge)
	}
    return nil
}

type Claims struct {
	UserID   uint32 `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RefreshClaims struct {
    UserID uint32 `json:"user_id"`
    Username string `json:"username"`
    Email string `json:"email"`
    Subject string `json:"sub"`
}

func (c *Claims) GetUserID() uint32 {
	return c.UserID
}

func (c *Claims) GetUsername() string {
	return c.Username
}

func (c *Claims) GetEmail() string {
	return c.Email
}

func (c *Claims) Validate() error {
    if c.Username == "" {
        return errors.New("Username is missing")
    } else if c.Email == "" {
        return errors.New("Email is missing")
    }
    return nil
}
