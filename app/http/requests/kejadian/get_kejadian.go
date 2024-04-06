package kejadian

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type GetKejadian struct {
	IdTypeKejadian string `form:"id_type_kejadian" json:"id_type_kejadian" binding:"required"`
}

func (r *GetKejadian) Authorize(ctx http.Context) error {
	return nil
}

func (r *GetKejadian) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id_type_kejadian": "required",
	}
}

func (r *GetKejadian) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"id_type_kejadian.required": "ID Type Kejadian Cannot be Empty",
	}
}

func (r *GetKejadian) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *GetKejadian) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
