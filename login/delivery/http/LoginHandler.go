package http

import (
	"os"
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
	grp.POST("/refreshToken", handler.RefreshToken)
	grp.POST("/operator", handler.Save)
	grp.GET("/operators", handler.GetAllOperators)
	grp.POST("/operator/:id", handler.Update)
	grp.DELETE("/operator/:id", handler.Delete)
}

func (loginHandler *LoginHandler) Authenticate(echoCtx echo.Context) (err error) {

	username := echoCtx.FormValue("username")
	useremail := echoCtx.FormValue("useremail")

	requestCtx := echoCtx.Request().Context()
	err = loginHandler.LoginUC.Authenticate(requestCtx, username, useremail)

	fmt.Printf("Username: %s \nUserEmail: %s \n", username, useremail)

	if err != nil {
		return echoCtx.JSON(http.StatusNotFound, "User not found in database. Please check your credentials.")
	}

	return echoCtx.JSON(http.StatusOK, generateTokenPair(username, useremail))

	// // Create token
	// token := jwt.New(jwt.SigningMethodHS256)

	// // Set claims
	// claims := token.Claims.(jwt.MapClaims)
	// claims["name"] = username
	// claims["email"] = useremail
	// claims["authorized"] = true
	// claims["app"] = "exam105"
	// claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

	// t, err := token.SignedString([]byte(os.Getenv("ENV_ACCESS_TOKEN_SECRET"))) 
	// if err != nil {
	// 	return err
	// }

	// // Generate Refresh-Token
	// refreshToken := jwt.New(jwt.SigningMethodHS256)
	// rtClaims := refreshToken.Claims.(jwt.MapClaims)
	// rtClaims["name"] = username
	// rtClaims["email"] = useremail
	// rtClaims["authorized"] = true
	// rtClaims["app"] = "exam105"
	// rtClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	// rt, err := refreshToken.SignedString([]byte(os.Getenv("ENV_REFRESH_TOKEN_SECRET")))
	// if err != nil {
	// 	return err
	// }

	// return echoCtx.JSON(http.StatusOK, map[string]string{
	// 	"access_token": t,
	// 	"refresh_token": rt,
	// })
	
}

func (loginHandler *LoginHandler) RefreshToken(echoCtx echo.Context) (err error) {
	
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}

	tokenReq := tokenReqBody{}
	echoCtx.Bind(&tokenReq)

	token, err := jwt.Parse(tokenReq.RefreshToken, func(token *jwt.Token) (interface{}, error) {

		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ENV_REFRESH_TOKEN_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		username := fmt.Sprintf("%v", claims["name"]) 
		useremail := fmt.Sprintf("%v", claims["email"])

		requestCtx := echoCtx.Request().Context()
		err = loginHandler.LoginUC.Authenticate(requestCtx, username, useremail)	
		if err != nil {
			return echoCtx.JSON(http.StatusNotFound, "User not found in database. Please check your credentials.")
		}
	
		if claims["app"] == "exam105" && claims["authorized"] == true {

			newTokenPair := generateTokenPair(username, useremail)
			// if err != nil {
			// 	return err
			// }

			return echoCtx.JSON(http.StatusOK, newTokenPair)
		}

		return echoCtx.JSON(http.StatusUnauthorized, "Refresh token issued")
	}

	return err
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

func generateTokenPair(username string, useremail string) (map[string]string) {

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = username
	claims["email"] = useremail
	claims["authorized"] = true
	claims["admin"] = false
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

	access_token, _ := token.SignedString([]byte(os.Getenv("ENV_ACCESS_TOKEN_SECRET")))
	// if err != nil {
	// 	return nil, err
	// }

	//Refresh Token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["name"] = username
	rtClaims["email"] = useremail
	rtClaims["app"] = "exam105"
	rtClaims["authorized"] = true
	rtClaims["exp"] = time.Now().Add(time.Hour * 3).Unix()

	refresh_token, _ := refreshToken.SignedString([]byte(os.Getenv("ENV_REFRESH_TOKEN_SECRET")))
	// if err != nil {
	// 	return nil, err
	// }

	return map[string]string{
		"access_token":  access_token,
		"refresh_token": refresh_token,
	}
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
