package http

import (
	"github.com/sirupsen/logrus"
	"net/http"
	// "os"
	// "github.com/labstack/echo/middleware"
	"github.com/labstack/echo"
	"github.com/exam105-UPD/backend/domain"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// QuestionHandler  service layer
type SearchHandler struct {
	SearchUC domain.SearchUsecase
}

// NewSearchHandler will initialize the resources endpoint
func NewSearchHandler(e *echo.Echo, searchUseCase domain.SearchUsecase) {

	handler := &SearchHandler{
		SearchUC: searchUseCase,
	}

	// Restricted group
	grp := e.Group("dashboard/de")
	// grp.Use(middleware.JWT([]byte(os.Getenv("ENV_ACCESS_TOKEN_SECRET")))) // The string "secret" should be accessed from data entry. For details, https://echo.labstack.com/cookbook/jwt

	grp.POST("/search/date", handler.SearchBySingleDate)
	grp.POST("/search/daterange", handler.SearchByDateRange)
}

func (searchHandler *SearchHandler) SearchBySingleDate(echoCtx echo.Context) (error) {

	// username, _ := restricted(echoCtx)
	var searchParameters domain.SearchParameterByDate
	err := echoCtx.Bind(&searchParameters)
	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	requestCtx := echoCtx.Request().Context()
	result, err := searchHandler.SearchUC.SearchByDate(requestCtx, &searchParameters) 
	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, result)

}

func (searchHandler *SearchHandler) SearchByDateRange(echoCtx echo.Context) (error) {

	// username, _ := restricted(echoCtx)
	var searchParameters domain.SearchParameterByDateRange
	err := echoCtx.Bind(&searchParameters)
	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	requestCtx := echoCtx.Request().Context()
	result, err := searchHandler.SearchUC.SearchByDateRange(requestCtx, &searchParameters) 
	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, result)

}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}