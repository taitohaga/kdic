package route

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/taitohaga/kdic/config"
	"github.com/taitohaga/kdic/services/auth"
)

func CreateAuthRoute(p iris.Party) {
	verifyMiddleware := config.Verifier.Verify(func() interface{} {
		return new(config.Claims)
	})
	p.Handle("GET", "/", servicePing)
	p.Handle("POST", "/login", getJWT)
	p.Handle("POST", "/refresh", refreshJWT)
	p.Handle("POST", "/create", createUser)
	p.Handle("GET", "/i/{username:string}", getUser)
	p.Handle("GET", "/i/{user_id:uint32}", getUserWithID)

	profile := p.Party("/profile")
	profile.Use(verifyMiddleware)
	profile.Handle("GET", "/", getProfile)
	profile.Handle("GET", "/email", getEmail)
}

func servicePing(ctx iris.Context) {
	ctx.JSON(iris.Map{"status": 1, "service": "Authentication of users"})
}

func getJWT(ctx iris.Context) {
	var loginReq auth.LoginRequest
	err := ctx.ReadJSON(&loginReq)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Login Failure").DetailErr(err))
		return
	}
	response, loginErr := auth.GetJWT(loginReq)
	if loginErr != nil {
		ctx.StatusCode(iris.StatusBadRequest)
	}
	ctx.JSON(response)
}

func refreshJWT(ctx iris.Context) {
	var request auth.RefreshJWTRequest
	err := ctx.ReadJSON(&request)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Refresh Failure").DetailErr(err))
		return
	}
	response, refreshJWTErr := auth.RefreshJWT(auth.RefreshJWTRequest{RefreshToken: request.RefreshToken})
	if refreshJWTErr != nil {
		ctx.StatusCode(iris.StatusBadRequest)
	}
	ctx.JSON(response)
}

func createUser(ctx iris.Context) {
	var newUser auth.CreateUserRequest
	err := ctx.ReadJSON(&newUser)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Failed to create user").DetailErr(err))
		return
	}
	response, createuserErr := auth.CreateUser(newUser)
	if createuserErr != nil {
		ctx.StatusCode(iris.StatusBadRequest)
	}
	ctx.JSON(response)
}

func getUser(ctx iris.Context) {
	username := ctx.Params().GetString("username")
	response, err := auth.GetUser(auth.GetUserRequest{UserName: username})
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
	}
	ctx.JSON(response)
}

func getUserWithID(ctx iris.Context) {
	userID, err := ctx.Params().GetUint32("user_id")
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Failed to fetch user").DetailErr(err))
		return
	}
	response, getErr := auth.GetUserWithID(auth.GetUserWithIDRequest{UserID: userID})
	if getErr != nil {
		ctx.StatusCode(iris.StatusNotFound)
	}
	ctx.JSON(response)
}

func getProfile(ctx iris.Context) {
	username := jwt.Get(ctx).(*config.Claims).Username
	response, err := auth.GetUser(auth.GetUserRequest{UserName: username})
	if err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().Title("Unauthorized").Detail(
			fmt.Sprintf("Token given but could not fetch your profile: %s", err),
		))
	}
	ctx.JSON(response.User)
}

func getEmail(ctx iris.Context) {
	username := jwt.Get(ctx).(*config.Claims).Username
	response, err := auth.GetEmail(auth.GetEmailRequest{Username: username})
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
	}
	ctx.JSON(response)
}
