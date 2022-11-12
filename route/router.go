package route

import (
    "github.com/kataras/iris/v12"
    "github.com/taitohaga/kdic/services/version"
)

type Route struct {
    Version *iris.Party
    Auth *iris.Party
}

func CreateRoute(app *iris.Application) Route {
    version := app.Party("version")
    auth := app.Party("auth")
    return Route{Version: &version, Auth: &auth}
}

func CreateVersionRoute(p iris.Party) {
    p.Handle("GET", "/", func (ctx iris.Context) {
        ctx.JSON(iris.Map{"version": version.GetVersion()})
    })
    p.Handle("GET", "/detail", func (ctx iris.Context) {
        ctx.JSON(version.GetDetailedInfo())
    })
}
