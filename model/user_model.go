package model

type User struct {
	UserUuid string `json:"userUuid"`
	Name     string `json:"name" validate:"required"`
	PhotoUrl string `json:"photoUrl"`
}

type CreateUserResponse struct {
	UserUuid string `json:"userUuid"`
}
