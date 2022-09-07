package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/timshannon/badgerhold"
	"github.com/txgruppi/study-server/dto"
	"github.com/txgruppi/study-server/models"
	"github.com/txgruppi/study-server/utils"
)

type PostCreateRequestData struct {
	Title string `json:"title" validate:"required"`
	Text  string `json:"text"`
}

type PostListQueryParams struct {
	Sort  *string `query:"sort" validate:"omitempty,oneof='title' '-title' 'created_at' '-created_at' 'updated_at' '-updated_at'"`
	Skip  *int    `query:"skip" validate:"omitempty,min=0"`
	Limit *int    `query:"limit" validate:"omitempty,min=1"`
}

type PostPatchRequestData struct {
	Title *string `json:"title"`
	Text  *string `json:"text"`
}

type Posts struct {
	Controller
	Store *badgerhold.Store
}

func (t *Posts) Create(c echo.Context) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	var reqData PostCreateRequestData
	if err := c.Bind(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	reqData.Title = strings.TrimSpace(reqData.Title)
	reqData.Text = strings.TrimSpace(reqData.Text)
	if err := c.Validate(&reqData); err != nil {
		return err
	}

	postID, err := utils.GenerateID()
	if err != nil {
		return err
	}

	now := time.Now()
	event := models.EventPostCreated{
		ID:        postID,
		Title:     reqData.Title,
		Text:      reqData.Text,
		CreatedAt: now,
		UpdatedAt: now,
	}
	doc := models.Post{}
	doc.AddAndApply(&event)

	if err := t.Store.Insert(postID, &doc); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, dto.NewPostFromModel(&doc))
}

func (t *Posts) List(c echo.Context) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	var params PostListQueryParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&params); err != nil {
		return err
	}

	query := &badgerhold.Query{}

	if params.Sort != nil {
		switch *params.Sort {
		case "title":
			{
				query = query.SortBy("Title")
			}

		case "-title":
			{
				query = query.SortBy("Title").Reverse()
			}

		case "updated_at":
			{
				query = query.SortBy("UpdatedAt")
			}

		case "-updated_at":
			{
				query = query.SortBy("UpdatedAt").Reverse()
			}

		case "created_at":
			{
				query = query.SortBy("CreatedAt")
			}

		case "-created_at":
			{
				query = query.SortBy("CreatedAt").Reverse()
			}
		}
	}
	if params.Skip != nil {
		query = query.Skip(*params.Skip)
	}
	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}

	var docs []*models.Post
	if err := t.Store.Find(&docs, query); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	res := make([]dto.Post, len(docs))
	for i, doc := range docs {
		res[i] = *dto.NewPostFromModel(doc)
	}

	return c.JSON(http.StatusOK, res)
}

func (t *Posts) Patch(c echo.Context) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	var reqData PostPatchRequestData
	if err := c.Bind(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&reqData); err != nil {
		return err
	}
	postID := c.Param("postID")
	var found *models.Post
	err := t.Store.UpdateMatching(&models.Post{}, badgerhold.Where(badgerhold.Key).Eq(postID), func(record interface{}) error {
		doc, ok := record.(*models.Post)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "record is not a Post")
		}
		found = doc

		changeTitle := reqData.Title != nil && *reqData.Title != doc.Title
		changeText := (reqData.Text != nil && *reqData.Text != doc.Text) || (reqData.Text == nil && doc.Text != "")
		text := ""
		if reqData.Text != nil {
			text = *reqData.Text
		}

		var event models.PostEvents
		if changeTitle && changeText {
			event = &models.EventPostTitleAndTextUpdated{
				Title:     *reqData.Title,
				Text:      text,
				UpdatedAt: time.Now(),
			}
		} else if reqData.Title != nil && *reqData.Title != doc.Title {
			event = &models.EventPostTitleUpdated{
				Title:     *reqData.Title,
				UpdatedAt: time.Now(),
			}
		} else if reqData.Text != nil && *reqData.Text != doc.Text {
			event = &models.EventPostTextUpdated{
				Text:      text,
				UpdatedAt: time.Now(),
			}
		}
		if event == nil {
			return nil
		}

		doc.AddAndApply(event)

		return nil
	})

	if err != nil {
		return err
	}

	if found == nil {
		return echo.NewHTTPError(http.StatusNotFound, "post not found")
	}

	found.ID = postID
	return c.JSON(http.StatusOK, dto.NewPostFromModel(found))
}

func (t *Posts) GetByID(c echo.Context) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	postID := c.Param("postID")
	var doc *models.Post
	err := t.Store.Get(postID, &doc)
	if err == badgerhold.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound, "post not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, dto.NewPostWithVersionsFromModel(doc))
}

func (t *Posts) DeleteByID(c echo.Context) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	postID := c.Param("postID")
	err := t.Store.Delete(postID, &models.Post{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
