package http

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/exam105-UPD/backend/domain"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// QuestionHandler  service layer
type QuestionHandler struct {
	QuestionUC domain.QuestionUsecase
}

// NewQuestionHandler will initialize the quesrion/ resources endpoint
func NewQuestionHandler(e *echo.Echo, qsUseCase domain.QuestionUsecase) {

	handler := &QuestionHandler{
		QuestionUC: qsUseCase,
	}

	// Restricted group
	grp := e.Group("dashboard/de")
	grp.Use(middleware.JWT([]byte("secret"))) // The string "secret" should be accessed from data entry. For details, https://echo.labstack.com/cookbook/jwt

	//e.GET("/articles", handler.FetchArticle)
	grp.GET("/test", handler.Testing)
	grp.POST("/questions", handler.Save)
	//e.GET("/articles/:id", handler.GetByID)
	//e.DELETE("/articles/:id", handler.Delete)
	grp.GET("/metadata/", handler.GetMetadataByUsernameAndEmail)
}

func (qsHandler *QuestionHandler) Testing(echoCtx echo.Context) (err error) {

	return echoCtx.JSON(http.StatusOK, "Get method called. Testing successful.")
}

func (qsHandler *QuestionHandler) Save(echoCtx echo.Context) (err error) {

	username, useremail := restricted(echoCtx)
	var allQuestion domain.MCQModel
	err = echoCtx.Bind(&allQuestion)
	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	// var ok bool
	// if ok, err = isRequestValid(&article); !ok {
	// 	return c.JSON(http.StatusBadRequest, err.Error())
	// }

	requestCtx := echoCtx.Request().Context()
	err = qsHandler.QuestionUC.Save(requestCtx, &allQuestion, username, useremail)
	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusCreated, allQuestion)
}

func (qsHandler *QuestionHandler) GetMetadataByUsernameAndEmail(echoCtx echo.Context) (err error) {

	username, useremail := restricted(echoCtx)
	requestCtx := echoCtx.Request().Context()

	metadataInfo, err := qsHandler.QuestionUC.GetMetadataById(requestCtx, username, useremail)
	
	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, metadataInfo)

}

func restricted(c echo.Context) (string, string) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	email := claims["email"].(string)

	return name, email
	// return c.String(http.StatusOK, "Welcome "+name+"!")
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
