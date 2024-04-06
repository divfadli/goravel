package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type RekapKejadianKecelakaanController struct {
	//Dependent services
}

func NewRekapKejadianKecelakaanController() *RekapKejadianKecelakaanController {
	return &RekapKejadianKecelakaanController{
		//Inject services
	}
}

func (r *RekapKejadianKecelakaanController) Index(ctx http.Context) http.Response {
	return nil
}	
