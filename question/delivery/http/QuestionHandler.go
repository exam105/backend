package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

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
	QuestionUC  domain.QuestionUsecase
	awsS3Client *s3.Client
}

// NewQuestionHandler will initialize the quesrion/ resources endpoint
func NewQuestionHandler(e *echo.Echo, qsUseCase domain.QuestionUsecase, thisAwsClient *s3.Client) {

	handler := &QuestionHandler{
		QuestionUC:  qsUseCase,
		awsS3Client: thisAwsClient,
	}

	// Restricted group
	grp := e.Group("dashboard/de")
	grp.Use(middleware.JWT([]byte(os.Getenv("ENV_ACCESS_TOKEN_SECRET")))) // The string "secret" should be accessed from data entry. For details, https://echo.labstack.com/cookbook/jwt

	grp.POST("/questions", handler.SaveMCQ)
	grp.GET("/metadata", handler.GetMetadataByUser)
	grp.POST("/metadata/:id", handler.UpdateMetadataByUser)
	grp.DELETE("/metadata/:id", handler.DeleteMetadataByUser)

	grp.GET("/questions/:metaid", handler.GetListOfMCQsByMetadataID)
	grp.GET("/question/:id", handler.GetQuestionByID)
	grp.POST("/question/:id", handler.UpdateQuestionByID)
	grp.PUT("/question/meta/:metaid", handler.AddQuestion)
	grp.DELETE("/question/:id/meta/:metaid", handler.DeleteQuestionByID)

	// Theory Questions
	grp.POST("/questions/theory", handler.SaveTheoryQs)
	grp.GET("/questions/theory/:metaid", handler.GetListOfMCQsByMetadataID)
	grp.GET("/question/theory/:id", handler.GetTheoryQuestionByID)
	grp.POST("/question/theory/:id", handler.UpdateTheoryQuestionByID)
	grp.PUT("/question/theory/meta/:metaid", handler.AddTheoryQuestion)
	grp.DELETE("/question/theory/:id/meta/:metaid", handler.DeleteQuestionByID)

	//S3 Credentials
	grp.GET("/question/s3credentials", handler.GetS3Credentials)

	//JWT Free URLs
	grp2 := e.Group("exam")
	grp2.GET("/test", handler.Testing)
	grp2.GET("/question/:id", handler.GetQuestionByID_NoAuth) // This will return Theory and MCQ question
	grp2.GET("/questions/theory/:metaid", handler.GetListOfMCQsByMetadataID_NoAuth)
	grp2.GET("/questions/:metaid", handler.GetListOfMCQsByMetadataID_NoAuth)
	grp2.GET("/metadata/:metaid", handler.GetMetadataById_NoAuth)
	grp2.GET("/env", handler.GetEnvVariables)
	grp2.POST("/question/uploadimage", handler.UploadImageToS3)
	grp2.GET("/homepage/links", handler.FetchHomePageLinks)
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

func (qsHandler *QuestionHandler) GetMetadataByUser(echoCtx echo.Context) error {

	username, useremail := restricted(echoCtx)
	requestCtx := echoCtx.Request().Context()

	metadataInfo, err := qsHandler.QuestionUC.GetMetadataById(requestCtx, username, useremail)

	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, metadataInfo)

}

func (qsHandler *QuestionHandler) UpdateMetadataByUser(echoCtx echo.Context) error {

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

func (qsHandler *QuestionHandler) DeleteMetadataByUser(echoCtx echo.Context) error {

	_, _ = restricted(echoCtx)

	docID := echoCtx.Param("id")
	requestCtx := echoCtx.Request().Context()
	metadataInfo, err := qsHandler.QuestionUC.DeleteMetadataById(requestCtx, docID)

	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, metadataInfo)
}

func (qsHandler *QuestionHandler) GetListOfMCQsByMetadataID(echoCtx echo.Context) error {

	_, _ = restricted(echoCtx)

	metadataID := echoCtx.Param("metaid")
	requestCtx := echoCtx.Request().Context()

	allQuestion, err := qsHandler.QuestionUC.GetMCQsByMetadataID(requestCtx, metadataID)

	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, allQuestion)

}

func (qsHandler *QuestionHandler) GetQuestionByID(echoCtx echo.Context) error {

	_, _ = restricted(echoCtx)

	questionID := echoCtx.Param("id")
	requestCtx := echoCtx.Request().Context()

	question, err := qsHandler.QuestionUC.GetQuestion(requestCtx, questionID)

	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, question)

}

func (qsHandler *QuestionHandler) GetTheoryQuestionByID(echoCtx echo.Context) error {

	_, _ = restricted(echoCtx)

	questionID := echoCtx.Param("id")
	requestCtx := echoCtx.Request().Context()

	question, err := qsHandler.QuestionUC.GetTheoryQuestion(requestCtx, questionID)

	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, question)

}

func (qsHandler *QuestionHandler) UpdateQuestionByID(echoCtx echo.Context) error {

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

func (qsHandler *QuestionHandler) UpdateTheoryQuestionByID(echoCtx echo.Context) error {

	_, _ = restricted(echoCtx)

	docID := echoCtx.Param("id")

	var updatedQuestion domain.TheoryQuestion
	err := echoCtx.Bind(&updatedQuestion)

	requestCtx := echoCtx.Request().Context()

	questionResult, err := qsHandler.QuestionUC.UpdateTheoryQuestion(requestCtx, updatedQuestion, docID)

	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, questionResult)
}

func (qsHandler *QuestionHandler) DeleteQuestionByID(echoCtx echo.Context) error {

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

func (qsHandler *QuestionHandler) AddQuestion(echoCtx echo.Context) error {

	_, _ = restricted(echoCtx)
	metaID := echoCtx.Param("metaid")

	var singleQuestion domain.Question
	err := echoCtx.Bind(&singleQuestion)

	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	requestCtx := echoCtx.Request().Context()
	result, err := qsHandler.QuestionUC.AddSingleQuestion(requestCtx, &singleQuestion, metaID)

	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusCreated, result)

}

func (qsHandler *QuestionHandler) AddTheoryQuestion(echoCtx echo.Context) error {

	_, _ = restricted(echoCtx)
	metaID := echoCtx.Param("metaid")

	var singleQuestion domain.TheoryQuestion
	err := echoCtx.Bind(&singleQuestion)

	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	requestCtx := echoCtx.Request().Context()
	result, err := qsHandler.QuestionUC.AddSingleTheoryQuestion(requestCtx, &singleQuestion, metaID)

	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusCreated, result)

}

func (qsHandler *QuestionHandler) GetS3Credentials(echoCtx echo.Context) error {

	type S3Cred struct {
		Username  string `json:"username" bson:"username"`
		Accesskey string `json:"accesskey" bson:"accesskey"`
		Secretkey string `json:"secretkey" bson:"secretkey"`
		Region    string `json:"region" bson:"region"`
	}

	s3cred := new(S3Cred)
	s3cred.Username = os.Getenv("ENV_S3_USERNAME")
	s3cred.Accesskey = os.Getenv("ENV_S3_ACCESS_KEY_ID")
	s3cred.Secretkey = os.Getenv("ENV_S3_SECRET_ACCESS_KEY")
	s3cred.Region = os.Getenv("ENV_S3_REGION")

	//log.Info("S3 Region: " + os.Getenv("ENV_S3_REGION"))

	return echoCtx.JSON(http.StatusOK, s3cred)
}

func (qsHandler *QuestionHandler) UploadImageToS3(echoCtx echo.Context) error {

	file, fileHeader, err := echoCtx.Request().FormFile("file")
	subject := echoCtx.Request().FormValue("subject")
	if err != nil {
		return err
	}

	// Get the fileName from Path
	// imageFile := "/home/muhammad/Pictures/Image-1.jpeg"
	imageFile := fileHeader.Filename

	// // Open the file from the file path
	// upFile, err := os.Open("/home/muhammad/Pictures/Image-1.jpeg")
	// if err != nil {
	// 	return fmt.Errorf("could not open local filepath : %v", err)
	// }
	// defer upFile.Close()

	// // Get the file info
	// upFileInfo, _ := upFile.Stat()
	// var fileSize int64 = upFileInfo.Size()
	// fileBuffer := make([]byte, fileSize)
	// upFile.Read(fileBuffer)

	// filename := header.Filename

	// defer file.Close()

	var sb strings.Builder
	sb.WriteString(subject + "/")
	sb.WriteString(imageFile)

	fmt.Println("String Builder: " + sb.String())

	uploader := manager.NewUploader(qsHandler.awsS3Client)
	uploadResult, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("exam105"),
		Key:    aws.String(sb.String()),
		Body:   file,
	})

	if err != nil {
		fmt.Printf("Error: %v  \n", err)
		return err
	}

	uploadLocation := uploadResult.Location
	fmt.Println("Location: " + uploadResult.Location)
	fmt.Printf("Image has been uploaded ->>> %v \t \n", imageFile)

	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusCreated, uploadLocation)

}

func (qsHandler *QuestionHandler) FetchHomePageLinks(echoCtx echo.Context) error {

	requestInput := &s3.GetObjectInput{
		Bucket: aws.String("exam105"),
		Key:    aws.String("homepage.json"),
	}

	result, err := qsHandler.awsS3Client.GetObject(context.TODO(), requestInput)
	if err != nil {
		fmt.Println(err)
	}
	defer result.Body.Close()

	body1, err := ioutil.ReadAll(result.Body)
	if err != nil {
		fmt.Println(err)
	}
	bodyString1 := fmt.Sprintf("%s", body1)

	// uploader := manager.NewUploader(qsHandler.awsS3Client)
	// uploadResult, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
	// 	Bucket: aws.String("exam105"),
	// 	Key:    aws.String(sb.String()),
	// 	Body:   file,
	// })

	// if err != nil {
	// 	return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	// }

	return echoCtx.JSON(http.StatusOK, bodyString1)

}

func (qsHandler *QuestionHandler) GetQuestionByID_NoAuth(echoCtx echo.Context) error {

	// _, _ = restricted(echoCtx)

	questionID := echoCtx.Param("id")
	requestCtx := echoCtx.Request().Context()

	question, err := qsHandler.QuestionUC.GetQuestionNoAuth(requestCtx, questionID)

	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, question)

}

func (qsHandler *QuestionHandler) GetListOfMCQsByMetadataID_NoAuth(echoCtx echo.Context) error {

	// _, _ = restricted(echoCtx)

	metadataID := echoCtx.Param("metaid")
	requestCtx := echoCtx.Request().Context()

	allQuestion, err := qsHandler.QuestionUC.GetMCQsByMetadataID(requestCtx, metadataID)

	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, allQuestion)

}

func (qsHandler *QuestionHandler) GetMetadataById_NoAuth(echoCtx echo.Context) error {

	// _, _ = restricted(echoCtx)

	metadataID := echoCtx.Param("metaid")
	requestCtx := echoCtx.Request().Context()

	allQuestion, err := qsHandler.QuestionUC.GetMetadataInfoByMetaIDNoAuth(requestCtx, metadataID)

	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, allQuestion)

}

func (qsHandler *QuestionHandler) GetEnvVariables(echoCtx echo.Context) error {

	type Env struct {
		GoogleAnalyticsMeasurementID string `json:"googleanalytics" bson:"googleanalytics"`
	}

	envVariables := new(Env)
	envVariables.GoogleAnalyticsMeasurementID = os.Getenv("ENV_GOOGLE_ANALYTICS_MEASUREMENT_ID")

	return echoCtx.JSON(http.StatusOK, envVariables)
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
