package approval

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type GetApproval struct {
	Nik string `form:"nik" json:"nik"`
}

func (r *GetApproval) Authorize(ctx http.Context) error {
	return nil
}

func (r *GetApproval) Rules(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *GetApproval) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *GetApproval) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *GetApproval) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
