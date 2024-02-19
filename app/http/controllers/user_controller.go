package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type UserController struct {
	//Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		//Inject services
	}
}

func (r *UserController) Show(ctx http.Context) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"Hello": "Goravel",
	})
}
func (r *UserController) Login(ctx http.Context) http.Response {
	var loginRequest requests
	// var loginRequest requests
	// sanitize := Sanitize(ctx, &loginRequest)
	// if sanitize != nil {
	// 	return sanitize
	// }
	// return ctx.Response().Success().Json(http.Json{
	// 	"Hello": "Goravel",
	// })
}
