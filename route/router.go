package route

import (
	"github.com/kataras/iris/v12"
	"github.com/taitohaga/kdic/services/version"
)

func Cors(ctx iris.Context) {
    ctx.Header("Access-Control-Allow-Origin", "http://localhost:3000")
    ctx.Header("Access-Control-Allow-Credentials", "true")

    if ctx.Method() == iris.MethodOptions {
        ctx.Header("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS")
        ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Content-Type, Authorization, Max-Age")
        ctx.Header("Access-Control-Max-Age", "86400")
        ctx.StatusCode(iris.StatusNoContent)
        return
    }
    ctx.Next()
}

func CreateVersionRoute(p iris.Party) {
    p.UseRouter(Cors)
    p.AllowMethods(iris.MethodOptions)
	p.Handle("GET", "/", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"version": version.GetVersion()})
	})
	p.Handle("GET", "/detail", func(ctx iris.Context) {
		ctx.JSON(version.GetDetailedInfo())
	})
}
