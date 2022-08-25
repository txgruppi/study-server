package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/timshannon/badgerhold"
	"github.com/txgruppi/tasks-server/models"
	"github.com/txgruppi/tasks-server/utils"
)

type TaskCreateRequestData struct {
	Title string    `json:"title" validate:"required"`
	Start time.Time `json:"start" validate:"required"`
}

type TaskPatchRequestData struct {
	End time.Time `json:"end" validate:"required"`
}

type Tasks struct {
	Controller
	Store *badgerhold.Store
}

func (t *Tasks) Create(c echo.Context) error {
	var reqData TaskCreateRequestData
	if err := c.Bind(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	reqData.Title = strings.TrimSpace(reqData.Title)
	if err := c.Validate(&reqData); err != nil {
		return err
	}
	taskID, err := utils.GenerateID()
	if err != nil {
		return err
	}
	doc := models.Task{
		Title: reqData.Title,
		Start: reqData.Start,
	}
	if err := t.Store.Insert(taskID, &doc); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, doc)
}

func (t *Tasks) List(c echo.Context) error {
	var tasks []models.Task
	if err := t.Store.Find(&tasks, nil); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if tasks == nil {
		tasks = []models.Task{}
	}
	return c.JSON(http.StatusOK, tasks)
}

func (t *Tasks) Patch(c echo.Context) error {
	var reqData TaskPatchRequestData
	if err := c.Bind(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&reqData); err != nil {
		return err
	}
	taskID := c.Param("taskID")
	var found *models.Task
	err := t.Store.UpdateMatching(&models.Task{}, badgerhold.Where(badgerhold.Key).Eq(taskID), func(record interface{}) error {
		doc, ok := record.(*models.Task)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "record is not a Task")
		}
		if doc.End != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "task already completed")
		}
		if doc.Start.After(reqData.End) {
			return echo.NewHTTPError(http.StatusBadRequest, "end must be after start")
		}
		found = doc
		doc.End = &reqData.End
		return nil
	})
	if err != nil {
		return err
	}
	if found == nil {
		return echo.NewHTTPError(http.StatusNotFound, "task not found")
	}
	found.ID = taskID
	return c.JSON(http.StatusOK, found)
}

func (t *Tasks) GetByID(c echo.Context) error {
	taskID := c.Param("taskID")
	var doc *models.Task
	err := t.Store.Get(taskID, &doc)
	if err == badgerhold.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound, "task not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, doc)
}

func (t *Tasks) DeleteByID(c echo.Context) error {
	taskID := c.Param("taskID")
	err := t.Store.Delete(taskID, &models.Task{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
