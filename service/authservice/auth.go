package authservice

import "gin-blog/models"

type Auth struct {
	Username string
	PassWord string
}

func (a *Auth) Check() (bool, error) {
	return models.CheckAuth(a.Username, a.PassWord)
}
