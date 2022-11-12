package route

import (
    "github.com/kataras/iris/v12"
    "github.com/taitohaga/kdic/services/auth"
)

func CreateAuthRoute(p iris.Party) {
    p.Handle("GET", "/", servicePing)
    p.Handle("POST", "/login", getJWT)
    p.Handle("POST", "/create", createUser)
    p.Handle("GET", "/i/{username:string}", getUser)
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
        ctx.StopWithProblem(iris.StatusNotFound, iris.NewProblem().Title("Failed to get user info").DetailErr(err))
        return
    }
    ctx.JSON(response)
}
