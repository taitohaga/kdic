package main

import (
	"os"

	"github.com/kataras/iris/v12"
	"github.com/taitohaga/kdic/config"
	"github.com/taitohaga/kdic/route"
)

func main() {
	app := iris.Default()

	configErr := config.InitConfig()
	app.Logger().Debugf("Mode: %s", os.Getenv("MODE"))
	app.Logger().Debugf("DB Path: %s", os.Getenv("DBPATH"))
	app.Logger().Debugf("Secret Key: %s", os.Getenv("SECRET"))
	if configErr != nil {
		app.Logger().Fatalf("Failed to load config: %s", configErr)
	}

	app.Logger().Debugf("Connecting to DB (path: %s", os.Getenv("DBPATH"))
	dbErr := config.InitDb()
	if dbErr != nil {
		app.Logger().Fatalf("Failed to connect to DB: %s", dbErr)
	}
	jwtErr := config.InitJWT()
	if jwtErr != nil {
		app.Logger().Fatalf("Failed to initialize JWT configs: %s", jwtErr)
	}

	version := app.Party("version")
	route.CreateVersionRoute(version)

	auth := app.Party("auth")
	route.CreateAuthRoute(auth)

	dic := app.Party("dic")
	route.CreateDicRoute(dic)

	app.Listen(":8080")
}
