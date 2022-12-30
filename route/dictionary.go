package route

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"

	"github.com/taitohaga/kdic/config"
	"github.com/taitohaga/kdic/services/dic"
	"github.com/taitohaga/kdic/services/word"
)

func CreateDicRoute(p iris.Party) {
	verifyMiddleware := config.Verifier.Verify(func() interface{} {
		return new(config.Claims)
	})
	p.Use(verifyMiddleware)
	p.Use(authenticate)

	p.Handle("POST", "/create", createDictionary)
	p.Handle("GET", "/i/{dicname:string}", getDictionary)
	p.Handle("GET", "/i/{dicname:string}/people", checkAuthority)
	p.Handle("POST", "/edit/{dicname:string}", setDictionary)

	word := p.Party("/{dicname:string}")
	word.Handle("POST", "/create", createWord)
	word.Handle("POST", "/edit/{word_id:uint}", setWord)
    word.Handle("DELETE", "/delete", delWord)
	word.Handle("GET", "/words", listWord)
}

func authenticate(ctx iris.Context) {
	dn := ctx.Params().GetStringDefault("dicname", "")
	if dn == "" {
		ctx.Next()
	} else {
		checkAuthorityResponse, err := dic.CheckAuthority(
			dic.CheckAuthorityRequest{
				DictionaryName: dn,
			},
		)
		if err != nil {
			ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().Title("Cannot access the dictionary").DetailErr(err))
			return
		}
		userID := jwt.Get(ctx).(*config.Claims).UserID
		username := jwt.Get(ctx).(*config.Claims).Username
		isPermitted := false
		for _, u := range checkAuthorityResponse.Users {
			if u.ID == userID {
				isPermitted = true
			}
		}
		if isPermitted {
			ctx.Next()
		} else {
			ctx.StopWithProblem(
				iris.StatusUnauthorized,
				iris.NewProblem().Title("No permission").Detail(fmt.Sprintf("User %s is not allowed to edit dictionary %s", username, dn)),
			)
		}
	}
}

func createDictionary(ctx iris.Context) {
	var newDictionary dic.CreateDictionaryRequest
	err := ctx.ReadJSON(&newDictionary)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Failed to create dictionary").DetailErr(err))
		return
	}
	claims, _ := jwt.Get(ctx).(*config.Claims)
	response, createdicErr := dic.CreateDictionary(newDictionary, claims)
	if createdicErr != nil {
		ctx.StatusCode(iris.StatusBadRequest)
	}
	ctx.JSON(response)
}

func getDictionary(ctx iris.Context) {
	response, err := dic.GetDictionary(dic.GetDictionaryRequest{
		DictionaryName: ctx.Params().GetString("dicname"),
	})
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
	}
	ctx.JSON(response)
}

func checkAuthority(ctx iris.Context) {
	response, err := dic.CheckAuthority(dic.CheckAuthorityRequest{DictionaryName: ctx.Params().GetString("dicname")})
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
	}
	ctx.JSON(response)
}

func setDictionary(ctx iris.Context) {
	var request interface{}
	ctx.ReadJSON(&request)
	response, err := dic.SetDictionary(request, ctx.Params().GetString("dicname"))
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
	}
	ctx.JSON(response)
}

func createWord(ctx iris.Context) {
	var request word.CreateWordRequest
	ctx.ReadJSON(&request)
	claims, _ := jwt.Get(ctx).(*config.Claims)
	response, err := word.CreateWord(request, claims, ctx.Params().GetString("dicname"))
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
	}
	ctx.JSON(response)
}

func setWord(ctx iris.Context) {
	var request interface{}
	ctx.ReadJSON(&request)
	claims, _ := jwt.Get(ctx).(*config.Claims)
	response, err := word.SetWord(request, claims, ctx.Params().GetString("dicname"), ctx.Params().GetUintDefault("word_id", 0))
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
	}
	ctx.JSON(response)
}

func listWord(ctx iris.Context) {
	var request interface{}
	ctx.ReadJSON(&request)
	response, err := word.ListWord(ctx.Params().GetString("dicname"))
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
	}
	ctx.JSON(response)
}

func delWord(ctx iris.Context) {
	var request word.DelWordRequest
	ctx.ReadJSON(&request)
	claims, _ := jwt.Get(ctx).(*config.Claims)
	response, err := word.DelWord(request, claims, ctx.Params().GetString("dicname"))
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
	}
	ctx.JSON(response)
}
