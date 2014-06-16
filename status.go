package nap

import (
	"net/http"
)

type Status interface {
	Code() int
	Message() string
}

type OK struct {
	Msg string
}

func (s OK) Code() int {
	return http.StatusOK
}
func (s OK) Message() string {
	return s.Msg
}

type Created struct {
	Msg string
}

func (s Created) Code() int {
	return http.StatusCreated
}
func (s Created) Message() string {
	return s.Msg
}

type NotFound struct {
	Msg string
}

func (s NotFound) Code() int {
	return http.StatusNotFound
}
func (s NotFound) Message() string {
	return s.Msg
}

type BadRequest struct {
	Msg string
}

func (s BadRequest) Code() int {
	return http.StatusBadRequest
}
func (s BadRequest) Message() string {
	return s.Msg
}

type MethodNotAllowed struct {
	Msg string
}

func (s MethodNotAllowed) Code() int {
	return http.StatusMethodNotAllowed
}
func (s MethodNotAllowed) Message() string {
	return s.Msg
}

type InternalError struct {
	Msg string
}

func (s InternalError) Code() int {
	return http.StatusInternalServerError
}
func (s InternalError) Message() string {
	return s.Msg
}
