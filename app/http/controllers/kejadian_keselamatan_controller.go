package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type RekapKejadianKeselamatanController struct {
	//Dependent services
}

func NewRekapKejadianKeselamatanController() *RekapKejadianKeselamatanController {
	return &RekapKejadianKeselamatanController{
		//Inject services
	}
}

func (r *RekapKejadianKeselamatanController) Index(ctx http.Context) http.Response {
	return nil
}	
