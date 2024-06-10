package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/goravel/framework/contracts/http"
)

// SuccessResponse general success response
type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// ErrorResponse general error response
type ErrorResponse struct {
	Message string `json:"message"`
}

// ErrorMessage general error message
type ErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return (fe.Field() + " wajib di isi")
	}
	return "Unknown error"
}

// Success response is successful
func Success(ctx http.Context, data any) http.Response {
	return ctx.Response().Success().Json(&SuccessResponse{
		Message: "success",
		Data:    data,
	})
}

// Error response error
func Error(ctx http.Context, code int, message string) http.Response {
	return ctx.Response().Json(code, &ErrorResponse{
		Message: message,
	})
}

// ErrorSystem responds to system errors
func ErrorSystem(ctx http.Context, message string) http.Response {
	return ctx.Response().Json(http.StatusInternalServerError, &ErrorResponse{
		Message: message,
	})
}

// Sanitize Methode Post disinfection request parameters
func SanitizePost(ctx http.Context, request http.FormRequest) http.Response {
	if errors, err := ctx.Request().ValidateRequest(request); err != nil {
		return Error(ctx, http.StatusUnprocessableEntity, err.Error())
	} else if errors != nil {
		return Error(ctx, http.StatusUnprocessableEntity, errors.One())
	}

	return nil
}

// Sanitize Methode Get disinfection request parameters
func SanitizeGet(ctx http.Context, err error) http.Response {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		out := make([]ErrorMessage, len(ve))
		for i, fe := range ve {
			out[i] = ErrorMessage{Field: fe.Field(), Message: getErrorMsg(fe)}
		}
		return ctx.Response().Json(http.StatusUnprocessableEntity, gin.H{"errorcode_": http.StatusUnprocessableEntity, "errormsg_": out})
	}

	return nil
}

func buildFileIdentificator(origFileName string) string {
	ext := filepath.Ext(origFileName)
	unixTime := time.Now().Unix()
	str := origFileName + strconv.Itoa(int(unixTime))
	data := []byte(str)
	hasher := md5.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)
	hashedString := hex.EncodeToString(hash)
	finalName := hashedString + ext
	return finalName
}
