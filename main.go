package main

import (
	"context"
	"fmt"
	"log"
	"time"

	_loginHandler "github.com/exam105-UPD/backend/login/delivery/http"
	_loginRepo "github.com/exam105-UPD/backend/login/repository"
	_loginUseCase "github.com/exam105-UPD/backend/login/usecase"

	_questionHandler "github.com/exam105-UPD/backend/question/delivery/http"
	_questionRepo "github.com/exam105-UPD/backend/question/repository"
	_questionUseCase "github.com/exam105-UPD/backend/question/usecase"

	_middleware "github.com/exam105-UPD/backend/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoConnClient *mongo.Client

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	/* 	dbHost := viper.GetString(`database.host`)
	   	dbPort := viper.GetString(`database.port`)
	   	dbUser := viper.GetString(`database.user`)
	   	dbPass := viper.GetString(`database.pass`)
	   	dbName := viper.GetString(`database.name`)
	   	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	   	val := url.Values{}
	   	val.Add("parseTime", "1")
	   	val.Add("loc", "Asia/Jakarta")
	   	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	   	dbConn, err := sql.Open(`mysql`, dsn) */
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dbConn := initializeMongoDatabase(ctx)

	defer func() {
		err := dbConn.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := _middleware.InitMiddleware()
	e.Use(middL.CORS)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, remote_ip:${remote_ip}, host:${host} \n",
	}))

	// **** Article wiring ****

	/* 	authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	   	ar := _articleRepo.NewMysqlArticleRepository(dbConn)

	   	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	   	au := _articleUcase.NewArticleUsecase(ar, authorRepo, timeoutContext)
	   	_articleHttpDelivery.NewArticleHandler(e, au)
	*/
	// *********************************

	// ****** Question Wiring ******
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	qsRepo := _questionRepo.NewQuestionRepository(dbConn)
	qsUC := _questionUseCase.NewQuestionUsecase(qsRepo, timeoutContext)
	_questionHandler.NewQuestionHandler(e, qsUC)

	//**********Login Wiring**************

	//timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	loginRepo := _loginRepo.NewLoginRepository(dbConn)
	loginUC := _loginUseCase.NewLoginUsecase(loginRepo, timeoutContext)
	_loginHandler.NewLoginHandler(e, loginUC)

	log.Fatal(e.Start(viper.GetString("server.address")))
}

func initializeMongoDatabase(ctx context.Context) *mongo.Client {

	// Set client options
	//clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017") // IMP-> Set Connection in ENV variable
	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")
	clientOptions = clientOptions.SetMaxPoolSize(100)                       //100 is default driver setting

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal("Couldn't connect to the database \n", err)
		fmt.Errorf(err.Error())
	} else {
		fmt.Println(" New MongoDB connection created ! ")
	}

	MongoConnClient = client
	return MongoConnClient

}

/* func initializeLogger() *logrus.Logger {
		var filename string = "/var/log/exam105.log"

	   	// Create the log file if doesn't exist. And append to it if it already exists.
	   	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	   	Formatter := new(logrus.TextFormatter)
	   	// You can change the Timestamp format. But you have to use the same date and time.
	   	// "2006-02-02 15:04:06" Works. If you change any digit, it won't work
	   	// ie "Mon Jan 2 15:04:05 MST 2006" is the reference time. You can't change it
	   	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	   	Formatter.FullTimestamp = true
	   	logrus.SetFormatter(Formatter)

	   	if err != nil {
	   		// Cannot open log file. Logging to stderr
	   		fmt.Println(err)
	   	} else {
	   		logrus.SetOutput(file)
	   	}
} */
