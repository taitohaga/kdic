package main

import (
	"os"

	"github.com/kataras/iris/v12"
	"github.com/taitohaga/kdic/config"
	"github.com/taitohaga/kdic/route"
)

func main() {
	app := iris.Default()
    
    config.InitConfig()
    app.Logger().Debugf("Mode: %s", os.Getenv("MODE"))
    app.Logger().Debugf("DB Path: %s", os.Getenv("DBPATH"))

    version := app.Party("version")
    route.CreateVersionRoute(version)

    auth := app.Party("auth")
    route.CreateAuthRoute(auth)
    
	app.Listen(":8080")
}
