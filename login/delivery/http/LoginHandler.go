package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/exam105-UPD/backend/domain"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
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
	//e.DELETE("/articles/:id", handler.Delete)
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
	claims["name"] = "Jon Snow"
	claims["email"] = useremail
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret")) // This should come from Env variable
	if err != nil {
		return err
	}

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
