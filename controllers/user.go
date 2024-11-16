package controllers

import (
	"go-crud-boilerplate/utils"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type (
	user struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

var (
	users = []*user{}
	seq   = 1
	lock  = sync.Mutex{}
)

func GetUsers(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	return utils.SendResponse(c, users, "Success", http.StatusOK)
}
