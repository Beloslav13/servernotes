package models

type Person struct {
	TgUserId uint   `json:"tg_user_id" validate:"required"`
	Username string `json:"username" validate:"required"`
}

type Category struct {
	PersonId uint   `json:"person_id" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

type Note struct {
	PersonId   uint   `json:"person_id" validate:"required"`
	CategoryId uint   `json:"category_id" validate:"required"`
	Name       string `json:"name" validate:"required"`
}

type Response struct {
	Result  bool        `json:"result"`
	Message string      `json:"message,omitempty"`
	Err     interface{} `json:"error"`
}