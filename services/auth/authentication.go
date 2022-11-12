package auth

import (
	"errors"
	"fmt"
)

type CreateUserRequest struct {
	Id          int64  `json:"id"`
	UserName    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Profile     string `json:"profile"`
}

type CreateUserResponse struct {
	Message  string `json:"msg"`
	UserName string `json:"username"`
}

func CreateUser(u CreateUserRequest) (CreateUserResponse, error) {
	return CreateUserResponse{Message: "created user", UserName: u.UserName}, nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"msg"`
	Token   string `json:"token"`
}

func GetJWT(loginReq LoginRequest) (LoginResponse, error) {
	if loginReq.Username == "taitohaga" && loginReq.Password == "password" {
		return LoginResponse{Message: fmt.Sprintf("Logged in as %s", loginReq.Username), Token: "token"}, nil
	}

	return LoginResponse{Message: "User does not exist or password is incorrect"}, errors.New("Login Failed")
}

type GetUserRequest struct {
	UserName string `json:"username"`
}

type GetUserResponse struct {
	Message string            `json:"msg"`
	User    CreateUserRequest `json:"user"`
}

func GetUser(request GetUserRequest) (GetUserResponse, error) {
    if request.UserName == "notfounduser" {
        return GetUserResponse{Message: fmt.Sprintf("%s not found", request.UserName)}, errors.New("User Not Found")
    }
	return GetUserResponse{
		Message: fmt.Sprintf("Found %s", request.UserName),
		User: CreateUserRequest{
			Id:          1,
			UserName:    request.UserName,
			DisplayName: "John Doe",
			Email:       "john@company.com",
			Profile:     "Nice person",
		},
	}, nil
}
