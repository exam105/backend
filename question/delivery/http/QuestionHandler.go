package http

import (
	"os"
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
	grp.Use(middleware.JWT([]byte(os.Getenv("ENV_ACCESS_TOKEN_SECRET")))) // The string "secret" should be accessed from data entry. For details, https://echo.labstack.com/cookbook/jwt

	grp.GET("/test", handler.Testing)
	grp.POST("/questions", handler.SaveMCQ)
	grp.GET("/metadata", handler.GetMetadataByUser)
	grp.POST("/metadata/:id", handler.UpdateMetadataByUser)
	grp.DELETE("/metadata/:id", handler.DeleteMetadataByUser)

	grp.GET("/questions/:id", handler.GetListOfMCQsByMetadataID)
	grp.GET("/question/:id", handler.GetQuestionByID)
	grp.POST("/question/:id", handler.UpdateQuestionByID)
	grp.PUT("/question/meta/:metaid", handler.AddQuestion)
	grp.DELETE("/question/:id/meta/:metaid", handler.DeleteQuestionByID)

	// Theory Questions
	grp.POST("/questions/theory", handler.SaveTheoryQs)
	grp.GET("/question/theory/:id", handler.GetTheoryQuestionByID)

	grp.GET("/question/s3credentials",handler.GetS3Credentials)
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

func (qsHandler *QuestionHandler) SaveTheoryQs(echoCtx echo.Context) (err error) {

	username, useremail := restricted(echoCtx)
	var allQuestion domain.TheoryModel
	err = echoCtx.Bind(&allQuestion)
	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	requestCtx := echoCtx.Request().Context()
	err = qsHandler.QuestionUC.SaveTheoryQuestion(requestCtx, &allQuestion, username, useremail)
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

func (qsHandler *QuestionHandler) GetQuestionByID(echoCtx echo.Context) (error){

	_, _ = restricted(echoCtx)

	questionID := echoCtx.Param("id")	
	requestCtx := echoCtx.Request().Context()

	question, err := qsHandler.QuestionUC.GetQuestion(requestCtx, questionID)
	
	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, question)

}

func (qsHandler *QuestionHandler) GetTheoryQuestionByID(echoCtx echo.Context) (error){

	_, _ = restricted(echoCtx)

	questionID := echoCtx.Param("id")	
	requestCtx := echoCtx.Request().Context()

	question, err := qsHandler.QuestionUC.GetTheoryQuestion(requestCtx, questionID)
	
	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, question)

}

func (qsHandler *QuestionHandler) UpdateQuestionByID(echoCtx echo.Context) (error){
	
	_, _ = restricted(echoCtx)

	docID := echoCtx.Param("id")

	var updatedQuestion domain.Question
	err := echoCtx.Bind(&updatedQuestion)

	requestCtx := echoCtx.Request().Context()

	questionResult, err := qsHandler.QuestionUC.UpdateQuestion(requestCtx, updatedQuestion, docID)
	
	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, questionResult)
}

func (qsHandler *QuestionHandler) DeleteQuestionByID(echoCtx echo.Context) (error) {

	_, _ = restricted(echoCtx)

	questionID := echoCtx.Param("id")
	metaID := echoCtx.Param("metaid")

	requestCtx := echoCtx.Request().Context()
	questionResult, err := qsHandler.QuestionUC.DeleteQuestion(requestCtx, metaID, questionID)
	
	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, questionResult)
	
}

func (qsHandler *QuestionHandler) AddQuestion(echoCtx echo.Context) (error) {

	_, _ = restricted(echoCtx)
	metaID := echoCtx.Param("metaid")

	var singleQuestion domain.Question
	err := echoCtx.Bind(&singleQuestion)

	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	requestCtx := echoCtx.Request().Context()
	result, err  := qsHandler.QuestionUC.AddSingleQuestion(requestCtx, &singleQuestion, metaID)
	
	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusCreated, result)

}

func (qsHandler *QuestionHandler) GetS3Credentials(echoCtx echo.Context) (error){

	type S3Cred struct {
		Username       	string 				`json:"username" bson:"username"`
		Accesskey 		string 				`json:"accesskey" bson:"accesskey"`
		Secretkey 		string 				`json:"secretkey" bson:"secretkey"`
		Region			string				`json:"region" bson:"region"`
	}
	
	s3cred := new(S3Cred)
	s3cred.Username = os.Getenv("ENV_S3_USERNAME")
	s3cred.Accesskey = os.Getenv("ENV_S3_ACCESS_KEY_ID")
	s3cred.Secretkey = os.Getenv("ENV_S3_SECRET_ACCESS_KEY")
	s3cred.Region = os.Getenv("ENV_S3_REGION")

	//log.Info("S3 Region: " + os.Getenv("ENV_S3_REGION"))

	return echoCtx.JSON(http.StatusOK, s3cred)
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
