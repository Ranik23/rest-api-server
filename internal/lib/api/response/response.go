package response

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"strings"
)


type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK = "OK"
	StatusError = "Error"
)


func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error: msg,
	}
}


func ValidationError(errors validator.ValidationErrors) Response {
	var Errors []string

	for _, err := range errors {

		switch err.ActualTag(){
			
		case "required":
			Errors = append(Errors, fmt.Sprintf("field %s is required", err.Field()))
		case "url":
			Errors = append(Errors, fmt.Sprintf("field %s is not a valid url", err.Field()))
		default:
			Errors = append(Errors, fmt.Sprintf("filed is not valid", err.Field()))
		}
	}

	return Response {
		Status: StatusError,
		Error: strings.Join(Errors, ", "),
	}

}




}

