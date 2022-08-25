package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/timshannon/badgerhold"
	"github.com/txgruppi/tasks-server/http/controllers"
)

type Validator struct {
	validate *validator.Validate
}

func (t *Validator) Validate(i interface{}) error {
	if err := t.validate.Struct(i); err != nil {
		switch err := err.(type) {
		case validator.ValidationErrors:
			res := []string{}
			for _, e := range err {
				res = append(res, e.Error())
			}
			return echo.NewHTTPError(http.StatusBadRequest, res)
		default:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}
	return nil
}

func Wire(store *badgerhold.Store) (*echo.Echo, error) {
	e := echo.New()
	e.Validator = &Validator{validate: validator.New()}

	tasks := controllers.Tasks{Store: store}

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/tasks", tasks.Create)
	e.GET("/tasks", tasks.List)
	e.PATCH("/tasks/:taskID", tasks.Patch)
	e.GET("/tasks/:taskID", tasks.GetByID)
	e.DELETE("/tasks/:taskID", tasks.DeleteByID)

	return e, nil
}
