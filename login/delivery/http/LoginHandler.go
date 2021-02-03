package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/exam105-UPD/backend/domain"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// QuestionHandler  service layer
type LoginHandler struct {
	LoginUC domain.LoginUsecase
}

// NewLoginHandler will initialize the login/ resources endpoint
func NewLoginHandler(e *echo.Echo, loginUseCase domain.LoginUsecase) {

	handler := &LoginHandler{
		LoginUC: loginUseCase,
	}

	grp := e.Group("superuser")

	grp.POST("/login", handler.Authenticate)
	grp.POST("/operator", handler.Save)
	grp.GET("/operators", handler.GetAllOperators)
	grp.POST("/operator/:id", handler.Update)
	grp.DELETE("/operator/:id", handler.Delete)
}

func (loginHandler *LoginHandler) Authenticate(echoCtx echo.Context) (err error) {

	username := echoCtx.FormValue("username")
	useremail := echoCtx.FormValue("useremail")

	fmt.Printf("Username: %s \nUserEmail: %s \n", username, useremail)

	// Throws unauthorized error
	if username != "jon" || useremail != "abc@efg.xyz" {
		return echo.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = username
	claims["email"] = useremail
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret")) // This should come from Env variable
	if err != nil {
		return err
	}

	//QuestionHandler. 
	return echoCtx.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
	
}

func (loginHandler *LoginHandler) Save(echoCtx echo.Context) (err error) {

	var dataEntryOperatorModel domain.DataEntryOperatorModel
	err = echoCtx.Bind(&dataEntryOperatorModel)

	requestCtx := echoCtx.Request().Context()
	err = loginHandler.LoginUC.Save(requestCtx, &dataEntryOperatorModel)
	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusCreated, dataEntryOperatorModel.Username+" : Account created")
}

func (loginHandler *LoginHandler) GetAllOperators(echoCtx echo.Context) (err error) {

	requestCtx := echoCtx.Request().Context()
	operatorList, err := loginHandler.LoginUC.GetAllOperators(requestCtx)

	if err != nil {
		return echoCtx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return echoCtx.JSON(http.StatusOK, operatorList)
}

func (loginHandler *LoginHandler) Update(echoCtx echo.Context) (err error) {

	docID := echoCtx.Param("id")
	objID, err := primitive.ObjectIDFromHex(docID)

	if err != nil {
		fmt.Println("ObjectIDFromHex ERROR", err)
	} else {
		fmt.Println("ObjectIDFromHex:", objID)
	}

	var deoModel domain.DataEntryOperatorModel
	err = echoCtx.Bind(&deoModel)

	requestCtx := echoCtx.Request().Context()
	updated, err := loginHandler.LoginUC.Update(requestCtx, &deoModel, objID)
	return echoCtx.JSON(http.StatusOK, updated)
}

func (loginHandler *LoginHandler) Delete(echoCtx echo.Context) error {

	docID := echoCtx.Param("id")
	objID, err := primitive.ObjectIDFromHex(docID)

	if err != nil {
		fmt.Println("ObjectIDFromHex ERROR", err)
	} else {
		fmt.Println("ObjectIDFromHex:", objID)
	}

	requestCtx := echoCtx.Request().Context()
	delete, err := loginHandler.LoginUC.Delete(requestCtx, objID)
	return echoCtx.JSON(http.StatusOK, delete)

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
