package dto

type UserData struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type UserCreated struct {
	Username             string `json:"username" validate:"required"`
	Name                 string `json:"name" validate:"required"`
	Role                 string `json:"role" validate:"required"`
	Status               string `json:"status"`
	Password             string `json:"password" validate:"required,password,min=8"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type UserUpdateData struct {
	ID     string `json:"-"`
	Name   string `json:"name" validate:"required"`
	Role   string `json:"role" validate:"required"`
	Status string `json:"status" validate:"required"`
}

type UserUpdatePassword struct {
	ID                   string `json:"-"`
	Password             string `json:"password" validate:"required,password,min=8"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}
