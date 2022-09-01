package controllers

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/txgruppi/study-server/utils"
)

type Controller struct {
	mtx sync.Mutex
}

func (c *Controller) GenerateID() (string, error) {
	id, err := utils.GenerateID()
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return id, nil
}
