package models

import (
	"errors"
	"net/http"
)

var (
	ErrBadRequest          = errors.New("Bad request")
	ErrNotFound            = errors.New("Your requested item is not found")
	ErrConflict            = errors.New("Your item has already exist")
	ErrUnauthorized        = errors.New("User not authorised or not found")
	ErrInternalServerError = errors.New("Internal Server Error")
)

func GetStatusCodePost(err error) int {
	if err == nil {
		return http.StatusCreated
	}

	switch err {
	case ErrBadRequest: // 400
		return http.StatusBadRequest
	case ErrNotFound:
		return http.StatusNotFound // 404
	case ErrConflict:
		return http.StatusConflict // 409
	case ErrUnauthorized:
		return http.StatusUnauthorized // 401
	default:
		return http.StatusInternalServerError // 500
	}
}

func GetStatusCodeGet(err error) int {
	if err == nil {
		return http.StatusOK
	}
	return http.StatusNotFound
}
type Error struct {
	Message string `json:"message"`
}