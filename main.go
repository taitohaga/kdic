package main

import (
	"github.com/kataras/iris/v12"
)

type LoginRequest struct {
    Username string `json:"username"`
    Email string `json:"email"`
    Password string `json:"password"`
}

func main() {
    app := iris.Default()
    app.Handle("GET", "/", func(ctx iris.Context) {
        ctx.JSON(iris.Map{"status": 1})
    })

    app.Handle("POST", "/login", func(ctx iris.Context) {
        var loginReq LoginRequest
        err := ctx.ReadJSON(&loginReq)
        if err != nil {
            ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Login Failure").DetailErr(err))
            return
        }
        if loginReq.Username == "taitohaga" && loginReq.Password == "password" {
            ctx.JSON(iris.Map{"status": 1, "access_token": "token"})
        } else {
            ctx.StatusCode(iris.StatusBadRequest)
        }
    })

    app.Listen(":8080")
}
