package helpers

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

type BadRequest struct {
	Status  bool     `json:"status"`
	Message string   `json:"message"`
	Data    struct{} `json:"data"`
}

type EmptyResponse struct{}

func Success(message string, data interface{}, errors interface{}) Response {
	res := Response{
		Status:  true,
		Message: message,
		Errors:  nil,
		Data:    data,
	}

	return res
}

type ApiError struct {
	Field string `json:"field"`
	Msg   string `json:"message"`
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	}
	return ""
}

func Error(message string, err error, data interface{}) Response {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ApiError, len(ve))
		for i, fe := range ve {
			// fmt.Println(fe.Namespace())
			// fmt.Println(fe.Field())
			// fmt.Println(fe.StructNamespace())
			// fmt.Println(fe.StructField())
			// fmt.Println(fe.Tag())
			// fmt.Println(fe.ActualTag())
			// fmt.Println(fe.Kind())
			// fmt.Println(fe.Type())
			// fmt.Println(fe.Value())
			// fmt.Println(fe.Param())
			// fmt.Println()
			out[i] = ApiError{fe.Field(), fe.Tag()}
		}

		res := Response{
			Status:  false,
			Message: message,
			Errors:  out,
			Data:    data,
		}

		return res
	}

	res := Response{
		Status:  false,
		Message: message,
		Errors:  err,
		Data:    data,
	}

	return res
}

func BuildBadRequest(status bool, message string, data struct{}) BadRequest {
	res := BadRequest{
		Status:  status,
		Message: message,
		Data:    data,
	}

	return res
}
