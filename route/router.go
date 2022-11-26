package route

import (
    "github.com/kataras/iris/v12"
    "github.com/taitohaga/kdic/services/version"
)

func CreateVersionRoute(p iris.Party) {
    p.Handle("GET", "/", func (ctx iris.Context) {
        ctx.JSON(iris.Map{"version": version.GetVersion()})
    })
    p.Handle("GET", "/detail", func (ctx iris.Context) {
        ctx.JSON(version.GetDetailedInfo())
    })
}
