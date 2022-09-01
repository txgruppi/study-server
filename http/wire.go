package http

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/timshannon/badgerhold"
	"github.com/txgruppi/parseargs-go"
	"github.com/txgruppi/study-server/errors"
	"github.com/txgruppi/study-server/http/controllers"
)

type Validator struct {
	validate *validator.Validate
}

func (t *Validator) Validate(i interface{}) error {
	if err := t.validate.Struct(i); err != nil {
		switch err := err.(type) {
		case validator.ValidationErrors:
			{
				res := []errors.ValidationError{}
				for _, e := range err {
					next := errors.ValidationError{
						Field:    e.Field(),
						Tag:      e.Tag(),
						Value:    e.Value(),
						ErrorStr: e.Error(),
						Param:    e.Param(),
					}
					switch next.Tag {
					case "oneof":
						{
							param, err0 := parseargs.Parse(e.Param())
							if err0 == nil {
								next.Param = param
							}
						}

					case "min":
						{
							param, err0 := strconv.ParseInt(e.Param(), 10, 64)
							if err0 == nil {
								next.Param = param
							}
						}

					case "required":
						{
							next.Param = nil
							next.Value = nil
						}
					}
					res = append(res, next)
				}
				return echo.NewHTTPError(http.StatusBadRequest, res)
			}

		case *errors.ValidationError:
			{
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

		default:
			{
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
		}
	}
	return nil
}

func Wire(store *badgerhold.Store) (*echo.Echo, error) {
	e := echo.New()
	e.Validator = &Validator{validate: validator.New()}

	tasks := controllers.Tasks{Store: store}
	posts := controllers.Posts{Store: store}

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/tasks", tasks.Create)
	e.GET("/tasks", tasks.List)
	e.PATCH("/tasks/:taskID", tasks.Patch)
	e.GET("/tasks/:taskID", tasks.GetByID)
	e.DELETE("/tasks/:taskID", tasks.DeleteByID)

	e.POST("/posts", posts.Create)
	e.GET("/posts", posts.List)
	e.PATCH("/posts/:postID", posts.Patch)
	e.GET("/posts/:postID", posts.GetByID)
	e.DELETE("/posts/:postID", posts.DeleteByID)

	return e, nil
}
