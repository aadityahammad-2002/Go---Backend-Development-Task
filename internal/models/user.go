package models

import "time"

type UserRequest struct {
	Name string `json:"name" validate:"required"`
	DOB  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

type UserCreateResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
}

type UserUpdateResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
}

type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Age  int    `json:"age"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type PaginatedUsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

type DBUser struct {
	ID   int32
	Name string
	DOB  time.Time
}
