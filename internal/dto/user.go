package dto

type CreateUserRequest struct {
	FirstName string `json:"first_name"binding:"required"`
	LastName  string `json:"last_name"binding:"required"`
	Email     string `json:"email"binding:"required""`
	Password  string `json:"password"binding:"required"`
	Phone     string `json:"phone"`
}

type UserResponse struct {
	UserId    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
type CreateUserResponse struct {
	AccessToken string       `json:"accessToken"`
	User        UserResponse `json:"user"`
}

type LoginRequest struct {
	Email    string `json:"email"binding:"required"`
	Password string `json:"password"binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	User        struct {
		UserId    string `json:"userId"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	} `json:"user"`
}
