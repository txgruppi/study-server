package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/txgruppi/tasks-server/utils"
)

type Controller struct {
}

func (c *Controller) GenerateID() (string, error) {
	id, err := utils.GenerateID()
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return id, nil
}
