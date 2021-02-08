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

	grp.GET("/test", handler.Testing)
	grp.POST("/questions", handler.SaveMCQ)
	grp.GET("/metadata", handler.GetMetadataByUser)
	grp.POST("/metadata/:id", handler.UpdateMetadataByUser)
	grp.DELETE("/metadata/:id", handler.DeleteMetadataByUser)

	grp.GET("/questions/:id", handler.GetListOfMCQsByMetadataID)
}

func (qsHandler *QuestionHandler) Testing(echoCtx echo.Context) (err error) {

	return echoCtx.JSON(http.StatusOK, "Get method called. Testing successful.")
}

func (qsHandler *QuestionHandler) SaveMCQ(echoCtx echo.Context) (err error) {

	username, useremail := restricted(echoCtx)
	var allQuestion domain.MCQModel
	err = echoCtx.Bind(&allQuestion)
	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	requestCtx := echoCtx.Request().Context()
	err = qsHandler.QuestionUC.SaveMCQ(requestCtx, &allQuestion, username, useremail)
	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusCreated, allQuestion)
}

func (qsHandler *QuestionHandler) GetMetadataByUser(echoCtx echo.Context) (error) {

	username, useremail := restricted(echoCtx)
	requestCtx := echoCtx.Request().Context()

	metadataInfo, err := qsHandler.QuestionUC.GetMetadataById(requestCtx, username, useremail)
	
	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, metadataInfo)

}

func (qsHandler *QuestionHandler) UpdateMetadataByUser(echoCtx echo.Context) (error) {

	_, _ = restricted(echoCtx)

	docID := echoCtx.Param("id")

	var receivedMetadata domain.MetadataBson
	err := echoCtx.Bind(&receivedMetadata)

	requestCtx := echoCtx.Request().Context()

	metadataInfo, err := qsHandler.QuestionUC.UpdateMetadataById(requestCtx, receivedMetadata, docID)
	
	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, metadataInfo)

}

func (qsHandler *QuestionHandler) DeleteMetadataByUser(echoCtx echo.Context) (error) {

	// Delete Metadata should be able to delete all the question related to the paper. 
	// This function is ONLY deleting the Metadata but the question remains in the database which is not the expected behaviour.

	_, _ = restricted(echoCtx)

	docID := echoCtx.Param("id")
	requestCtx := echoCtx.Request().Context()
	metadataInfo, err := qsHandler.QuestionUC.DeleteMetadataById(requestCtx, docID)
	
	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, metadataInfo)
}

func (qsHandler *QuestionHandler) GetListOfMCQsByMetadataID(echoCtx echo.Context) (error){

	_, _ = restricted(echoCtx)

	metadataID := echoCtx.Param("id")	
	requestCtx := echoCtx.Request().Context()

	allQuestion, err := qsHandler.QuestionUC.GetMCQsByMetadataID(requestCtx, metadataID)
	
	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, allQuestion)

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
