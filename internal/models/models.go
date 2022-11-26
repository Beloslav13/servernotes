package models

import "time"

type Person struct {
	Id       uint   `json:"id"`
	TgUserId uint   `json:"tg_user_id" validate:"required"`
	Username string `json:"username" validate:"required"`
}

type Category struct {
	Id       uint   `json:"id"`
	PersonId uint   `json:"person_id" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

type Note struct {
	Id         uint      `json:"id"`
	PersonId   uint      `json:"person_id" validate:"required"`
	CategoryId uint      `json:"category_id" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	Created    time.Time `json:"created"`
}

type response struct {
	Result  bool          `json:"result"`
	Message string        `json:"message"`
	Err     interface{}   `json:"error"`
	Obj     []interface{} `json:"obj,omitempty"`
}

func NewResponse(result bool, message string, err, obj interface{}) response {
	return response{
		Result:  result,
		Message: message,
		Err:     err,
		Obj:     []interface{}{obj},
	}
}
