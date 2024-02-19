package controllers

import (
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
		Message: "Error: " + message,
	})
}

// ErrorSystem responds to system errors
func ErrorSystem(ctx http.Context) http.Response {
	return ctx.Response().Json(http.StatusInternalServerError, &ErrorResponse{
		Message: "System internal error",
	})
}

// Sanitize disinfection request parameters
func Sanitize(ctx http.Context, request http.FormRequest) http.Response {
	errors, err := ctx.Request().ValidateRequest(request)
	if err != nil {
		return Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if errors != nil {
		return Error(ctx, http.StatusUnprocessableEntity, errors.One())
	}

	return nil
}
