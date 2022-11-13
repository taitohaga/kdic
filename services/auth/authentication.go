package auth

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/taitohaga/kdic/config"
	"github.com/taitohaga/kdic/model"
	"github.com/taitohaga/kdic/util"
	"gorm.io/gorm"
)

var usernamePattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]{3,299}$`)
var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9_+-]+(.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\.)+[a-zA-Z]{2,}$`)
var passwordPattern = regexp.MustCompile(`^[a-zA-Z0-9.?/-_!*:;'"\^&@#$+=]{8,300}`)

type CreateUserRequest struct {
	Id          uint32 `json:"id"`
	UserName    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Profile     string `json:"profile"`
}

type CreateUserResponse struct {
	Message  string `json:"msg"`
	UserName string `json:"username"`
	UserID   uint32 `json:"user_id"`
}

// Create new user, set email and password.
func CreateUser(req CreateUserRequest) (CreateUserResponse, error) {
	u := model.User{
		Username:    req.UserName,
		DisplayName: req.DisplayName,
		Profile:     req.Profile,
	}

	{
		if !usernamePattern.MatchString(req.UserName) {
			return CreateUserResponse{
                Message:  fmt.Sprintf("Empty or invalid username: \"%s\"", req.UserName),
				UserName: req.UserName,
			}, errors.New("Empty or invalid username")
		}
		if !emailPattern.MatchString(req.Email) {
			return CreateUserResponse{
				Message:  fmt.Sprintf("Empty or invalid email: \"%s\"", req.Email),
				UserName: req.UserName,
			}, errors.New("Empty or invalid email")
		}
		if !passwordPattern.MatchString(req.Password) {
			return CreateUserResponse{
				Message:  "Empty or invalid password",
				UserName: req.UserName,
			}, errors.New("Empty or invalid password")
		}
	}

    err := config.Db.Transaction(func(tx *gorm.DB) error {
		insertUser := tx.Create(&u)
		if insertUser.Error != nil {
			return errors.New(fmt.Sprintf("failed to create user: %s", insertUser.Error))
		}
		u_email := model.Email{
			User:  u,
			Email: req.Email,
            IsPrimary: true,
		}
		insertEmail := tx.Create(&u_email)
		if insertEmail.Error != nil {
			return errors.New(fmt.Sprintf("Could not use given email: %s", insertEmail.Error))
		}
		u_account := model.Account{
			User:     u,
			Password: util.HashSHA256(req.Password),
		}
		insertAccount:= tx.Create(&u_account)
		if insertAccount.Error != nil {
			return errors.New(fmt.Sprintf("Could not set password: %s", insertAccount.Error))
		}
        return nil
	})
    if err != nil {
        return CreateUserResponse{
            Message: fmt.Sprintf("Failed to create user: %s", err),
        }, err
    }
	return CreateUserResponse{Message: "Created user", UserName: req.UserName, UserID: u.ID}, nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"msg"`
	AccessToken   string `json:"token"`
    RefreshToken string `json:"refresh_token"`
}

// Create new JSON Web Token from the given user information.
func GetJWT(request LoginRequest) (LoginResponse, error) {
    var u model.User
    var email model.Email
    var account model.Account
    if request.Username != "" {
        if res := config.Db.Where("username = ?", request.Username).First(&u); res.Error != nil {
            if errors.Is(res.Error, gorm.ErrRecordNotFound) {
                return LoginResponse{
                    Message: fmt.Sprintf("User not fonud: %s", request.Username),
                }, res.Error
            } else {
                return LoginResponse{
                    Message: fmt.Sprintf("Failed to fetch user info: %s", res.Error),
                }, res.Error
            }
        }
        if res := config.Db.Where("user_id = ? AND is_primary = ?", u.ID, true).First(&email); res.Error != nil {
            if errors.Is(res.Error, gorm.ErrRecordNotFound) {
                return LoginResponse{
                    Message: "The user does not have primary email!",
                }, res.Error
            } else {
                return LoginResponse{
                    Message: fmt.Sprintf("Could not retrieve email from username: %s", res.Error),
                }, res.Error
            }
        }
    }

    if res := config.Db.Where("user_id = ?", u.ID).First(&account); res.Error != nil {
        return LoginResponse{
            Message: "User is unavailable!",
        }, res.Error
    }
    if account.Password != util.HashSHA256(request.Password) {
        return LoginResponse{
            Message: "User not found or incorrect password",
        }, errors.New("User not found or incorrect password")
    }

    accessClaims := config.Claims{
        UserID: u.ID,
        Username: u.Username,
        Email: email.Email,
    }
    refreshClaims := config.RefreshClaims{
        UserID: u.ID,
        Username: u.Username,
        Email: email.Email,
        Subject: u.Username,
    }
    tokenPair, err := config.Signer.NewTokenPair(accessClaims, refreshClaims, config.RefreshTokenMaxAge)
    if err != nil {
        return LoginResponse{
            Message: fmt.Sprintf("Could not generate token: %s", err),
        }, err
    }
    return LoginResponse{
        Message: fmt.Sprintf("Logged in as %s", u.Username),
        AccessToken: strings.Replace(string(tokenPair.AccessToken), `"`, ``, 2),
        RefreshToken: strings.Replace(string(tokenPair.RefreshToken), `"`, ``, 2),
    }, nil
}

type RefreshJWTRequest struct {
    RefreshToken string `json:"refresh_token"`
}

func RefreshJWT(request RefreshJWTRequest) (LoginResponse, error) {
    refreshToken := []byte(request.RefreshToken)
    result, err := config.Verifier.VerifyToken(refreshToken)
    if err != nil {
        return LoginResponse{
            Message: fmt.Sprintf("Could not refresh token: %s", err),
        }, err
    }
    var rc config.RefreshClaims
    result.Claims(&rc)
    getUserRes, getUserErr := GetUser(GetUserRequest{UserName: rc.Username})
    if getUserErr != nil {
        return LoginResponse{
            Message: fmt.Sprintf("Refresh token is invalid since username has been changed."),
        }, err
    }
    accessClaims := config.Claims{
        UserID: getUserRes.User.ID,
        Username: getUserRes.User.Username,
        Email: rc.Email,
    }
    refreshClaims := config.RefreshClaims{
        UserID: getUserRes.User.ID,
        Username: getUserRes.User.Username,
        Email: rc.Email,
        Subject: rc.Subject,
    }
    tokenPair, err := config.Signer.NewTokenPair(accessClaims, refreshClaims, config.RefreshTokenMaxAge)
    if err != nil {
        return LoginResponse{
            Message: fmt.Sprintf("Could not generate token: %s", err),
        }, err
    }
    return LoginResponse{
        Message: fmt.Sprintf("Logged in as %s", getUserRes.User.Username),
        AccessToken: strings.Replace(string(tokenPair.AccessToken), `"`, ``, 2),
        RefreshToken: strings.Replace(string(tokenPair.RefreshToken), `"`, ``, 2),
    }, nil
}

type GetUserRequest struct {
	UserName string `json:"username"`
}

type GetUserResponse struct {
	Message string            `json:"msg"`
	User    model.User `json:"user"`
}

// Get user information from username.
func GetUser(request GetUserRequest) (GetUserResponse, error) {
    u := model.User{Username: request.UserName}
    selectUser := config.Db.Where("username = ?", request.UserName).First(&u)
    if errors.Is(selectUser.Error, gorm.ErrRecordNotFound) {
        return GetUserResponse{
            Message: fmt.Sprintf("User not found: %s", request.UserName),
        }, selectUser.Error
    }
    return GetUserResponse{
        Message: fmt.Sprintf("Found user %s", u.Username),
        User: u,
    }, nil
}

type GetUserWithIDRequest struct {
    UserID uint32 `json:"user_id"`
}

// Get user information from user id.
func GetUserWithID(request GetUserWithIDRequest) (GetUserResponse, error) {
    u := model.User{ID: request.UserID}
    selectUser := config.Db.First(&u)
    if errors.Is(selectUser.Error, gorm.ErrRecordNotFound) {
        return GetUserResponse{
            Message: fmt.Sprintf("User not found: %d", request.UserID),
        }, selectUser.Error
    }
    return GetUserResponse{
        Message: fmt.Sprintf("Found user %d (%s)", u.ID, u.Username),
        User: u,
    }, nil
}

type GetEmailRequest struct {
    Username string `json:"username"`
}

type GetEmailResponse struct {
    Message string `json:"msg"`
    Email []model.Email `json:"email"`
}

func GetEmail(request GetEmailRequest) (GetEmailResponse, error) {
    var email_list []model.Email
    if res := config.Db.Joins("User").Find(&email_list); res.Error != nil {
        return GetEmailResponse{
            Message: fmt.Sprintf("Could not list user email: %s", res.Error),
        }, res.Error
    }
    return GetEmailResponse{
        Message: fmt.Sprintf("%s's %d email addresses fetched", request.Username, len(email_list)),
        Email: email_list,
    }, nil
}
