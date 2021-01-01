package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/exam105-UPD/backend/domain"
	"github.com/labstack/echo"
)

// QuestionHandler  service layer
type LoginHandler struct {
	LoginUC domain.LoginUsecase
}

// NewLoginHandler will initialize the login/ resources endpoint
func NewLoginHandler(e *echo.Echo, loginUseCase domain.LoginUsecase) {

	handler := &LoginHandler{
		LoginUC: loginUseCase,
	}

	e.POST("/login", handler.Authenticate)
	//e.POST("/questions", handler.Save)
	//e.GET("/articles/:id", handler.GetByID)
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
