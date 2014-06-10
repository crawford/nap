package nap

import (
	"net/http"
)

type Status interface {
	Code() int
	Message() string
}

type OK struct {
	message string
}

func (s OK) Code() int {
	return http.StatusOK
}
func (s OK) Message() string {
	return s.message
}

type Created struct {
	message string
}

func (s Created) Code() int {
	return http.StatusCreated
}
func (s Created) Message() string {
	return s.message
}

type NotFound struct {
	message string
}

func (s NotFound) Code() int {
	return http.StatusNotFound
}
func (s NotFound) Message() string {
	return s.message
}

type BadRequest struct {
	message string
}

func (s BadRequest) Code() int {
	return http.StatusBadRequest
}
func (s BadRequest) Message() string {
	return s.message
}

type MethodNotAllowed struct {
	message string
}

func (s MethodNotAllowed) Code() int {
	return http.StatusMethodNotAllowed
}
func (s MethodNotAllowed) Message() string {
	return s.message
}

type InternalError struct {
	message string
}

func (s InternalError) Code() int {
	return http.StatusInternalServerError
}
func (s InternalError) Message() string {
	return s.message
}
