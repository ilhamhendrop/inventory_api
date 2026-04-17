package dto

import "net/http"

type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func ResponseError(message string) Response[string] {
	return Response[string]{
		Code:    http.StatusInternalServerError,
		Message: message,
		Data:    "",
	}
}

func ResponseAuthError(message string) Response[string] {
	return Response[string]{
		Code:    http.StatusUnauthorized,
		Message: message,
		Data:    "",
	}
}

func ResponseDataError(message string, data map[string]string) Response[map[string]string] {
	return Response[map[string]string]{
		Code:    http.StatusBadRequest,
		Message: message,
		Data:    data,
	}
}

func ResponseSucces[T any](data T) Response[T] {
	return Response[T]{
		Code:    http.StatusOK,
		Message: "Succes",
		Data:    data,
	}
}

func ResponseCreated[T any](data T) Response[T] {
	return Response[T]{
		Code:    http.StatusCreated,
		Message: "Created",
		Data:    data,
	}
}

func ResponseNoContent() Response[string] {
	return Response[string]{
		Code:    http.StatusNoContent,
		Message: "No Content",
		Data:    "",
	}
}
